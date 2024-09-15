# Makefile for Go project

# Targets
.PHONY: all

all: build

build: vet test
	@echo "Compiling the project..."
	mkdir -p ./build
	go build -o ./build/ ./... 

test:
	go test ./...

run: build
	@echo "Running the application..."
	./build/echolocation

vet:
	@echo "Static analysis..."
	go vet ./...

mod-tidy:
	@echo "Tidying up modules..."
	go mod tidy

# Clean target
clean:
	@echo "Cleaning up..."
	go clean
	rm -rf ./build