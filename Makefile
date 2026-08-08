.PHONY: all test coverage bench lint clean e2e install-tools

# Default target
all: test

# Run all tests
test:
	go test -race ./...

# Run tests with coverage
coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Run linter
lint:
	golangci-lint run ./...

# Run end-to-end tests
e2e:
	go test -v ./e2e/...

# Run E2E tests with Serpro integration
e2e-serpro:
	go test -v ./e2e/... -run TestSerproIntegration

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	rm -rf dist/ build/

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build (for verification)
build:
	go build ./...
