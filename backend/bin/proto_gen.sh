#!/bin/bash

# Stop when error
set -e

# Set proto file path
PROTO_DIR="../api/proto"
# Set target path
GENERATED_DIR="../api/generated/ocr"

# Check existence of gen directory
mkdir -p $GENERATED_DIR

# Compile
protoc \
  --proto_path=$PROTO_DIR \
  --go_out=$GENERATED_DIR --go_opt=paths=source_relative \
  --go-grpc_out=$GENERATED_DIR --go-grpc_opt=paths=source_relative \
  $PROTO_DIR/ocr_service.proto
#protoc \
#  --proto_path=$PROTO_DIR \
#  --go_out=$GENERATED_DIR \
#  --go-grpc_out=$GENERATED_DIR \
#  $PROTO_DIR/ocr_service.proto

# Check result
if [ $? -eq 0 ]; then
  echo "Proto files successfully compiled to $GENERATED_DIR"
else
  echo "Failed to compile proto files"
  exit 1
fi
