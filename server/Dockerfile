FROM golang:1.17 as builder
WORKDIR /workspace
COPY . /workspace
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /workspace/server .
ENTRYPOINT ["/server"]
