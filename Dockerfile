# syntax=docker/dockerfile:1
FROM golang:1.21-alpine
LABEL authors="Elsa Lau"

WORKDIR /app

COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY db/example.db ./db/tasks.db

RUN go mod download
RUN go mod tidy
RUN go build -o server ./cmd/tasks

EXPOSE 8080

ENTRYPOINT ["./server"]
