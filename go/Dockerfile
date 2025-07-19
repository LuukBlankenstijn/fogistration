FROM golang:alpine3.22

RUN go install github.com/air-verse/air@latest && \
    go install github.com/jackc/tern/v2@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT ["sleep","3"]
