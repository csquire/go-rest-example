package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"

	"github.com/csquire/go-rest-example/application/handlers"
	"github.com/csquire/go-rest-example/infrastructure/persistence"
	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	log     = logrus.New()
	cfgFile string
)

const defaultPort = "8080"

func main() {
	pflag.StringVar(&cfgFile, "config", "config.yaml", "Config file path")
	pflag.Parse()

	viper.AutomaticEnv()

	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Warn(err)
	}

	basicAuthUser := viper.GetString("basicAuthUser")
	basicAuthPass := viper.GetString("basicAuthPass")

	db, err := persistence.NewPostgresDb(
		viper.GetString("databaseHost"),
		viper.GetString("databasePort"),
		viper.GetString("databaseUser"),
		viper.GetString("databasePassword"),
		viper.GetString("databaseName"),
		viper.GetString("databaseSslMode"))

	if err != nil {
		log.Fatalf("Unable to create database client: %v", err)
	}
	defer db.Close()

	router := newRouter(db, basicAuthUser, basicAuthPass)

	port := viper.GetString("serverPort")
	if port == "" {
		log.Printf("Default server to port %s", defaultPort)
		port = defaultPort
	}

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", port),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func newRouter(db *sqlx.DB, user, password string) http.Handler {
	repo := persistence.NewPostgresMetadataRepository(db)
	api := handlers.NewMetadataHandlers(repo)
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthCheck(db).ReadyEndpoint)
	metadataRouter := r.PathPrefix("/api/docker/metadata").Subrouter()
	metadataRouter.Use(httpauth.SimpleBasicAuth(user, password))
	metadataRouter.HandleFunc("", api.CreateMetadata).Methods("POST")
	metadataRouter.HandleFunc("", api.GetAllMetadata).Methods("GET")
	metadataRouter.HandleFunc("/{id}", api.GetMetadata).Methods("GET")
	metadataRouter.HandleFunc("/{id}/approve", api.Approve).Methods("POST")
	metadataRouter.HandleFunc("/{id}/deny", api.Deny).Methods("POST")

	return r
}

func healthCheck(db *sqlx.DB) healthcheck.Handler {
	health := healthcheck.NewHandler()
	health.AddReadinessCheck("db-readiness", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		res, err := db.ExecContext(ctx, "SELECT 1")
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows != 1 {
			return fmt.Errorf("incorrect DB result (expected 1 row, got %d)", rows)
		}
		return nil
	})
	return health
}
