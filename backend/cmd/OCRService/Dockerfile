# 构建阶段
FROM golang:1.23 AS builder

WORKDIR /app

# 复制依赖文件并安装构建所需的依赖
COPY go.mod go.sum ./
RUN go mod download

# 安装构建所需依赖
RUN apt-get update && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# 复制项目代码并编译
COPY cmd/OCRService ./cmd/OCRService
COPY internal ./internal
COPY pkg ./pkg
COPY api ./api

RUN GOOS=linux go build -o /app/OCRService ./cmd/OCRService

# Tesseract 构建阶段
FROM ubuntu:22.04 AS tesseract_builder

WORKDIR /build

# 安装构建 Tesseract 所需的依赖
RUN apt-get update && apt-get install -y \
    automake \
    autoconf \
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
    && rm -rf /var/lib/apt/lists/*

# 复制 Tesseract 源码
COPY engine/tesseract /build/tesseract

WORKDIR /build/tesseract

RUN apt-get update && apt-get install -y make

# 使用 autogen.sh 构建和安装 Tesseract
RUN chmod +x ./autogen.sh \
    && ./autogen.sh \
    && ./configure --prefix=/usr/local \
    && make -j$(nproc) \
    && make install \
    && ldconfig

# 运行时镜像
FROM ubuntu:22.04

WORKDIR /app

# 安装运行时依赖
RUN apt-get update && apt-get install -y \
    libleptonica-dev \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

# 从 Tesseract 构建阶段复制安装内容
COPY --from=tesseract_builder /usr/local /usr/local

# 复制应用程序二进制文件和配置
COPY --from=builder /app/OCRService /app/OCRService
COPY .env /app/.env

RUN ldconfig

# 暴露端口
EXPOSE 50051

# 设置默认 shell
SHELL ["/bin/bash", "-c"]

# 启动服务
CMD ["/app/OCRService"]
