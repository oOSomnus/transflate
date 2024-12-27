#!/bin/bash

# Stop when error
set -e

# Set proto file path and target path
PROTO_DIR="api/proto"
GENERATED_DIR="api/generated"

# Proto services to be compiled
SERVICES=("ocr" "translate")

# Ensure target directories exist
mkdir -p $GENERATED_DIR
for service in "${SERVICES[@]}"; do
  mkdir -p "$GENERATED_DIR/$service"
done

# Compile each service
for service in "${SERVICES[@]}"; do
  protoc \
    --proto_path=$PROTO_DIR \
    --go_out="$GENERATED_DIR/$service" --go_opt=paths=source_relative \
    --go-grpc_out="$GENERATED_DIR/$service" --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/${service}_service.proto"
done

echo "Proto files successfully compiled to $GENERATED_DIR"
