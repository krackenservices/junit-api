# Define binary name and the list of target operating systems and architectures
BINARY_NAME = junitapi

# Target platforms: darwin (macOS), linux, and windows for both amd64 and arm64 architectures
TARGETS = \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64 \
	windows/amd64 \
	windows/arm64

# Generate output directories and file extensions based on OS
OUTPUT_DIR = dist
WINDOWS_EXT = .exe

# Build modes
BUILD_MODE ?= debug

run:
	@go run *.go
# Default action: run tests and build binaries
all: test build

# Build binaries for all platforms
build: $(TARGETS)
$(TARGETS):
	@GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) \
		OUT_FILE=$(OUTPUT_DIR)/$(BUILD_MODE)/$(BINARY_NAME)-$(word 1,$(subst /, ,$@))-$(word 2,$(subst /, ,$@))$(if $(findstring windows,$(word 1,$(subst /, ,$@))),$(WINDOWS_EXT),) \
		&& echo "Building for OS: $(word 1,$(subst /, ,$@)), ARCH: $(word 2,$(subst /, ,$@)), Mode: $(BUILD_MODE) -> $$OUT_FILE" \
		&& mkdir -p $(OUTPUT_DIR)/$(BUILD_MODE) \
		&& CGO_ENABLED=0 GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) go build $(if $(filter release,$(BUILD_MODE)),-ldflags="-s -w",) -o $$OUT_FILE .

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)

# Build in debug mode
debug: BUILD_MODE=debug
debug: build

# Build in release mode
release: BUILD_MODE=release
release: build

.PHONY: all build test clean $(TARGETS) debug release
