package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/unrolled/render"

	"github.com/csquire/go-rest-example/domain/metadata"

	"github.com/gorilla/mux"
)

const ID = "id"

func NewMetadataHandlers(repo metadata.Repository) *MetadataApi {
	return &MetadataApi{
		render: render.New(),
		repo:   repo,
	}
}

// MetadataApi RPC API use for managing Docker metadata
// Errors are currently returned as plain text, not formatted as json.
type MetadataApi struct {
	render *render.Render
	repo   metadata.Repository
}

func (m *MetadataApi) CreateMetadata(w http.ResponseWriter, r *http.Request) {
	imageMetadata := metadata.ImageMetadata{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&imageMetadata)
	if err != nil {
		m.renderError(w, "Problem decoding request", err, http.StatusInternalServerError)
		return
	}

	if imageMetadata.Id == "" {
		m.renderError(w, "The id parameter cannot be blank", nil, http.StatusBadRequest)
		return
	}
	err = m.repo.CreateImage(&imageMetadata)

	if err != nil {
		m.renderError(w, "Problem creating metadata record", err, http.StatusInternalServerError)
		return
	}
	m.sendJsonResponse(w, http.StatusCreated, imageMetadata)
}

func (m *MetadataApi) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[ID]
	if id == "" {
		m.renderError(w, "The id parameter cannot be blank", nil, http.StatusBadRequest)
		return
	}

	imageMetadata, err := m.repo.GetImage(id)
	if err != nil {
		m.renderError(w, "Problem retrieving metadata", err, http.StatusInternalServerError)
		return
	}
	m.sendJsonResponse(w, http.StatusOK, imageMetadata)
}

func (m *MetadataApi) GetAllMetadata(w http.ResponseWriter, r *http.Request) {
	allImages, err := m.repo.GetAllImages()
	if err != nil {
		m.renderError(w, "Problem retrieving metadata", err, http.StatusInternalServerError)
		return
	}
	m.sendJsonResponse(w, http.StatusOK, allImages)
}

func (m *MetadataApi) Approve(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[ID]
	imageMetadata, err := m.approveOrDeny(id, true)
	if err != nil {
		m.renderError(w, "Problem approving metadata", err, http.StatusInternalServerError)
		return
	}
	m.sendJsonResponse(w, http.StatusOK, imageMetadata)
}

func (m *MetadataApi) Deny(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[ID]
	imageMetadata, err := m.approveOrDeny(id, false)
	if err != nil {
		m.renderError(w, "Problem denying metadata", err, http.StatusInternalServerError)
		return
	}
	m.sendJsonResponse(w, http.StatusOK, imageMetadata)
}

func (m *MetadataApi) approveOrDeny(id string, approved bool) (*metadata.ImageMetadata, error) {
	imageMetadata, err := m.repo.GetImage(id)
	if err != nil {
		return nil, err
	}
	imageMetadata.Approved = approved
	err = m.repo.UpdateImage(imageMetadata)
	if err != nil {
		return nil, err
	}
	return imageMetadata, nil
}

func (m *MetadataApi) sendJsonResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	err := m.render.JSON(w, statusCode, body)
	if err != nil {
		log.Printf("Problem rendering json response: %v", err)
	}
}

func (m *MetadataApi) renderError(w http.ResponseWriter, context string, err error, statusCode int) {
	msg := fmt.Sprintf("%s: %v", context, err)
	log.Print(msg)
	err = m.render.Text(w, statusCode, msg)
	if err != nil {
		log.Printf("Problem rendering response %v", err)
	}
}
