# basic image
FROM golang:1.23 AS builder

WORKDIR /app

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy related code
COPY cmd/TaskManager ./cmd/TaskManager
COPY internal ./internal
COPY pkg ./pkg
COPY api ./api


# compile TaskManager
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/TaskManager ./cmd/TaskManager

# runtime image
FROM alpine:3.18

WORKDIR /app

# install bash
RUN apk add --no-cache bash

# copy bin file built
COPY --from=builder /app/TaskManager /app/TaskManager
COPY .env /app/.env



EXPOSE 8080

SHELL ["/bin/bash", "-c"]

# server start
CMD ["/app/TaskManager"]
