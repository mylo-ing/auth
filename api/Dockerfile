#######################################################################
#  ── Dev stage ───────────────────────────────────────────────────── #
#  Hot-reload inside Docker (optional for local work)                 #
#######################################################################
FROM golang:1.24-alpine AS dev
WORKDIR /app

# Tools you wanted for Air + cgo + psql
RUN apk add --no-cache bash git build-base postgresql-client

# Speed up rebuilds
ENV GOPROXY=https://proxy.golang.org,direct \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/go-build

COPY go.mod go.sum ./
RUN go mod download

COPY api/ ./api

WORKDIR /app/api
RUN go install github.com/cosmtrek/air@v1.51.0
CMD ["air"]

#######################################################################
#  ── Prod build stage ────────────────────────────────────────────── #
#######################################################################
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 1) module files then deps
COPY go.mod go.sum ./
RUN go mod download

# 2) copy API code
COPY api/ ./api

# 3) compile the server (main package lives in /app/api)
WORKDIR /app/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-api .

#######################################################################
#  ── Tiny runtime stage ──────────────────────────────────────────── #
#######################################################################
FROM alpine:3.21 AS runtime
WORKDIR /app
COPY --from=builder /go-api /app/go-api

# expose the application port
EXPOSE 3517
CMD ["/app/go-api"]
