build:
	go build -o bin/server main.go

run:
	go run main.go

test:
	go test ./... -v -cover

docker:
	docker build -t csquire/go-rest-example:local .