# variables
PROTO_SRC_DIR := api/proto
PROTO_GEN_DIR := api/generated
PROTO_GEN_SCRIPT := bin/proto_gen.sh
GO_FILES := $(shell find . -type f -name '*.go')
OUTPUT_DIR := build

# services
SERVICES := OCRService TranslateService TaskManager

# default target
all: proto build

# gen proto file
proto:
	bash $(PROTO_GEN_SCRIPT)

# compile all services
build: proto
	@mkdir -p $(OUTPUT_DIR)
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build -o $(OUTPUT_DIR)/$$service ./cmd/$$service; \
	done

# clean
clean:
	rm -rf $(PROTO_GEN_DIR) $(OUTPUT_DIR)

# format
fmt:
	go fmt ./...

# static check
lint:
	golangci-lint run ./...

# test
test:
	go test ./...

# help
help:
	@echo "available targers："
	@echo "  all      - gen all proto files and compile"
	@echo "  proto    - gen proto files"
	@echo "  build    - compile all services"
	@echo "  clean    - clean gened files"
	@echo "  fmt      - format"
	@echo "  lint     - static check"
	@echo "  test     - run all tests"
	@echo "  help     - help info"

.PHONY: all proto build clean fmt lint test help
