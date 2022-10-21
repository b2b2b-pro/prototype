# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /http2repo ./cmd/app

## Deploy

FROM golang:1.18-alpine

WORKDIR /

COPY --from=builder /http2repo /http2repo

ENTRYPOINT ["/http2repo"]

EXPOSE 8088

