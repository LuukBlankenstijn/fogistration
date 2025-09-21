FROM node:24-bookworm-slim AS fe
WORKDIR /app
COPY frontend/.yarnrc.yml frontend/.yarn frontend/package.json frontend/yarn.lock ./
RUN corepack enable && yarn install --immutable
COPY /frontend ./
COPY /go/api /go/api
RUN yarn gen-client && yarn build

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache git
WORKDIR /src

# Cache Go modules & build cache
COPY go/go.mod go/go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./go .
COPY --from=fe /app/dist ./internal/http-server/http/spa/frontend/dist

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/cmdhandler ./cmd/cmdhandler && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/grpc        ./cmd/grpc        && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags docker_prod -o /out/http-server ./cmd/http-server


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
