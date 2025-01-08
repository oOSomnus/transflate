# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy codes
COPY cmd/task_manager ./cmd/task_manager
COPY internal/task_manager ./internal/task_manager
COPY pkg ./pkg
COPY api ./api

# compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/task_manager ./cmd/task_manager

FROM alpine:3.18

WORKDIR /app

# wkhtmltopdf
RUN apk add --no-cache bash curl

# copy compiled bin
COPY --from=builder /app/task_manager /app/task_manager
#COPY .env /app/.env
COPY config.local.yaml /app/config.local.yaml
COPY config.production.yaml /app/config.production.yaml
# expose
EXPOSE 8080

# run
CMD ["/app/task_manager"]
