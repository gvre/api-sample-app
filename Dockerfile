# syntax=docker/dockerfile:1

# Multi-stage build
# See https://docs.docker.com/develop/develop-images/multistage-build/ for details

# base
FROM golang:1.22-alpine AS golang
RUN apk add --update --no-cache alpine-sdk

# Copy go.mod and go.sum separately from the rest of the code,
# so their cached layer is not invalidated when the code changes.
COPY go.mod go.sum /
RUN go mod download

COPY . /app
WORKDIR /app/cmd/rest

# builder
FROM golang AS builder
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -o /rest -ldflags '-extldflags "-static" -s -w' .

# tester
FROM golang AS tester

# app
FROM alpine:latest

ENV PGHOST $PGHOST
ENV PGUSER $PGUSER
ENV PGPASSWORD $PGPASSWORD
ENV PGDATABASE $PGDATABASE

WORKDIR /
COPY --from=builder /rest .

EXPOSE 8080
CMD ["./rest"]
