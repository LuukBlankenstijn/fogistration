FROM golang:alpine3.22 AS dev

RUN go install github.com/air-verse/air@latest && \
    go install github.com/jackc/tern/v2@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    apk add --no-cache \
    ttf-dejavu \
    font-noto \
    font-noto-cjk \
    font-noto-emoji

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT ["sleep","3"]

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache git
WORKDIR /src

# Cache Go modules & build cache
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Cross-compile for the target platform (no QEMU)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/cmdhandler ./cmd/cmdhandler && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/grpc        ./cmd/grpc        && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/http-server ./cmd/http-server


FROM alpine:3.22 AS cmdhandler

WORKDIR /app
COPY --from=builder /out/cmdhandler /app/cmdhandler
ENTRYPOINT ["/app/cmdhandler"]


FROM alpine:3.22 AS grpc

RUN apk add --no-cache \
    ttf-dejavu \
    font-noto \
    font-noto-cjk \
    font-noto-emoji

WORKDIR /app
COPY --from=builder /out/grpc /app/grpc
EXPOSE 80
ENTRYPOINT ["/app/grpc"]


FROM alpine:3.22 AS http-server

WORKDIR /app
COPY --from=builder /out/http-server /app/http-server
EXPOSE 80
ENTRYPOINT ["/app/http-server"]
