BINARY_NAME=nectar
BUILD_PATH=bin/$(BINARY_NAME)
MAIN_PATH=$(BINARY_NAME)/main.go

VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build flags for setting version and build date
LDFLAGS := -ldflags "-X nectar/build.Version=$(VERSION) -X nectar/build.Date=$(BUILD_DATE)"

.PHONY: setup clean build run dev all

setup:
	@echo "Setting up environment..."
	@go mod download
	@echo "Environment setup complete."

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@echo "Cleanup complete."

build:
	@echo "Building version $(VERSION) at $(BUILD_DATE)..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o $(BUILD_PATH) $(MAIN_PATH) || true
	@echo "Build complete."

run:
	@if [ ! -f $(BUILD_PATH) ]; then echo "Binary not found. Building binary..."; $(MAKE) -s build; fi
	@echo "Running..."
	@$(BUILD_PATH) || true

dev:
	@echo "Running in development mode with version $(VERSION)..."
	@go run $(LDFLAGS) $(MAIN_PATH) || true

all: setup clean build run

.SILENT: