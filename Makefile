.PHONY: build run test proto clean help

# Build the server
build:
	@echo "Building server..."
	@go build -o bin/server cmd/server/main.go
	@echo "✅ Build complete: bin/server"

# Run the server
run:
	@echo "Starting Job Search Service..."
	@./bin/server

# Run test client
test:
	@echo "Running test client..."
	@go run cmd/client/test_client.go

# Generate proto files
proto:
	@echo "Generating proto files..."
	@export PATH=$$PATH:~/go/bin && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/job.proto
	@echo "✅ Proto files generated"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "✅ Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

# Check Elasticsearch
check-es:
	@echo "Checking Elasticsearch..."
	@curl -s http://localhost:9200 | jq .

# Build and run
start: build run

# Help
help:
	@echo "Job Search Service - Makefile Commands"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@echo "Commands:"
	@echo "  build      - Build the server binary"
	@echo "  run        - Run the server"
	@echo "  test       - Run the test client"
	@echo "  proto      - Regenerate proto files"
	@echo "  clean      - Remove build artifacts"
	@echo "  deps       - Install/update dependencies"
	@echo "  check-es   - Check Elasticsearch status"
	@echo "  start      - Build and run the server"
	@echo "  help       - Show this help message"
