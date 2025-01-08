# builder
FROM golang:1.23 AS builder

WORKDIR /app

# copy go dependencies
COPY go.mod go.sum ./
RUN go mod download

# install building dependencies
RUN apt-get update && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# copy code & compile
COPY cmd/ocr_service ./cmd/ocr_service
COPY internal ./internal
COPY pkg ./pkg
COPY api ./api

RUN GOOS=linux go build -o /app/ocr_service ./cmd/ocr_service

# Tesseract build
FROM ubuntu:22.04 AS tesseract_builder

WORKDIR /build

# install dependencies
RUN apt-get update && apt-get install -y \
    automake \
    autoconf \
    autoconf-archive \
    libtool \
    g++ \
    pkg-config \
    libpng-dev \
    libjpeg-dev \
    libtiff-dev \
    libicu-dev \
    zlib1g-dev \
    libpango1.0-dev \
    libcairo2-dev \
    libleptonica-dev \
    git \
    make \
    && rm -rf /var/lib/apt/lists/*

# copy origin code
#COPY engine/tesseract /build/tesseract
RUN git clone https://github.com/tesseract-ocr/tesseract.git /build/tesseract

WORKDIR /build/tesseract


RUN mkdir -p /build/tesseract/m4

# autogen install
RUN chmod +x ./autogen.sh \
    && ./autogen.sh \
    && ./configure --prefix=/usr/local \
    && make -j$(nproc) \
    && make install \
    && ldconfig

# runtime image
FROM ubuntu:22.04

WORKDIR /app

# install runtime dependencies
RUN apt-get update && apt-get install -y \
    libleptonica-dev \
    libgomp1 \
    poppler-utils \
    && rm -rf /var/lib/apt/lists/*

# merge previous stages
COPY --from=tesseract_builder /usr/local /usr/local
COPY --from=builder /app/ocr_service /app/ocr_service
#COPY .env /app/.env
COPY config.local.yaml /app/config.local.yaml
COPY config.production.yaml /app/config.production.yaml

COPY engine/ocr_trained_data/* /usr/local/share/tessdata

RUN ldconfig

# port
EXPOSE 50051

# default shell
SHELL ["/bin/bash", "-c"]

# start the service
CMD ["/app/ocr_service"]
