# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /repository ./cmd/app

## Deploy

FROM golang:1.18-alpine

WORKDIR /

COPY --from=builder /repository /repository

EXPOSE 50051

ENTRYPOINT ["/repository"]

