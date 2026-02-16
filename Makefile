BINARY   := image-processor
BUILD_DIR := ./bin

.PHONY: build run test lint clean

## build: compile the binary into ./bin/
build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/image-processor

## run: build and run with default sample images
run: build
	$(BUILD_DIR)/$(BINARY)

## test: run all tests with the race detector
test:
	go test -race -count=1 ./...

## lint: run go vet on all packages
lint:
	go vet ./...

## clean: remove build artifacts and processed output images
clean:
	rm -rf $(BUILD_DIR)
	rm -rf ./data/images/output/*