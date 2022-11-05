# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /auth ./cmd/app

## Deploy

FROM golang:1.18-alpine

WORKDIR /

COPY --from=builder /auth /auth

ENTRYPOINT ["/auth"]

EXPOSE 8888

