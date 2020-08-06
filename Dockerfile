# Multi-stage build
# See https://docs.docker.com/develop/develop-images/multistage-build/ for details

# builder
FROM golang:1.14-alpine AS builder

# Copy go.mod and go.sum separately from the rest of the code,
# so their cached layer is not invalidated when the code changes.
COPY go.mod go.sum /
RUN go mod download

COPY . /app
WORKDIR /app/cmd/rest
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -a -o /rest -ldflags '-extldflags "-static" -s -w' .

# app
FROM alpine:latest

ENV PGHOST $PGHOST
ENV PGUSER $PGUSER
ENV PGPASSWORD $PGPASSWORD
ENV PGDATABASE $PGDATABASE

WORKDIR /
COPY --from=builder /rest .

CMD ["./rest"]
