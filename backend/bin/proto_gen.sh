#!/bin/bash

# Stop when error
set -e

# Set proto file path
PROTO_DIR="api/proto"
# Set target path
GENERATED_DIR="api/generated"

# Check existence of gen directory
mkdir -p $GENERATED_DIR
mkdir -p $GENERATED_DIR/ocr
mkdir -p $GENERATED_DIR/translate

# Compile
protoc \
  --proto_path=$PROTO_DIR \
  --go_out=$GENERATED_DIR/ocr --go_opt=paths=source_relative \
  --go-grpc_out=$GENERATED_DIR/ocr --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/ocr_service.proto

protoc \
  --proto_path=$PROTO_DIR \
  --go_out=$GENERATED_DIR/translate --go_opt=paths=source_relative \
  --go-grpc_out=$GENERATED_DIR/translate --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/translate_service.proto


# Check result
if [ $? -eq 0 ]; then
  echo "Proto files successfully compiled to $GENERATED_DIR"
else
  echo "Failed to compile proto files"
  exit 1
fi
