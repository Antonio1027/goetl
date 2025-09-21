
# Start from the official Golang image
FROM golang:1.25.0-alpine AS builder
WORKDIR /app
ENV GO111MODULE=on
COPY . .
RUN apk update && apk upgrade --no-cache && go mod tidy && go build -o goetl ./cmd/main.go

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app
RUN apk update && apk upgrade --no-cache
COPY --from=builder /app/goetl /app/goetl
EXPOSE 8080
# Use sh -c to export env vars from .env before running the app
ENTRYPOINT ["sh", "-c", "exec ./goetl"]
