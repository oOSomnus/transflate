# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy codes
COPY cmd/translate_service ./cmd/translate_service
COPY internal/translate_service ./internal/translate_service
COPY pkg ./pkg
COPY api ./api

# compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/translate_service ./cmd/translate_service

FROM alpine:3.18

WORKDIR /app

# wkhtmltopdf
RUN apk add --no-cache bash curl

# copy compiled bin
COPY --from=builder /app/translate_service /app/translate_service
#COPY .env /app/.env
COPY config.local.yaml /app/config.local.yaml
COPY config.production.yaml /app/config.production.yaml
# expose
EXPOSE 50052

# run
CMD ["/app/translate_service"]
