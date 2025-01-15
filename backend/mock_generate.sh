#!/bin/bash

# 检查是否安装了 mockgen
if ! command -v mockgen &> /dev/null; then
    echo "Error: mockgen is not installed. Please install it using 'go install github.com/golang/mock/mockgen@latest'."
    exit 1
fi

# 检查参数
if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <interface_file>"
    echo "Example: $0 ./service.go"
    exit 1
fi

INTERFACE_FILE=$1

# 获取源文件所在目录和文件名
SOURCE_DIR=$(dirname "$INTERFACE_FILE")
SOURCE_FILE=$(basename "$INTERFACE_FILE")
OUTPUT_FILE="$SOURCE_DIR/mock_${SOURCE_FILE}"

# 提取包名
PACKAGE_NAME=$(grep -m 1 "^package" "$INTERFACE_FILE" | awk '{print $2}')

if [ -z "$PACKAGE_NAME" ]; then
    echo "Error: Unable to determine package name from $INTERFACE_FILE."
    exit 1
fi

# 提取所有接口名称
INTERFACES=$(grep -E '^type [A-Za-z0-9_]+ interface {' "$INTERFACE_FILE" | awk '{print $2}')

if [ -z "$INTERFACES" ]; then
    echo "No interfaces found in $INTERFACE_FILE."
    exit 1
fi

# 清空或创建目标文件
> "$OUTPUT_FILE"

# 运行 mockgen
for INTERFACE in $INTERFACES; do
    echo "Generating mock for interface: $INTERFACE"
    mockgen -source="$INTERFACE_FILE" -destination="$OUTPUT_FILE" -package="$PACKAGE_NAME" "$PACKAGE_NAME" "$INTERFACE"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to generate mock for interface $INTERFACE."
        exit 1
    fi
done

echo "All mocks generated successfully in: $OUTPUT_FILE"
