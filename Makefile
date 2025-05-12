# Define variables
APP_NAME = real-talk-forum
DB_FILE = reel-talk.db
SCHEMA_FILE = schema.sql

# Default target (runs the application)
run:
	go run main.go

# Build the application binary
build:
	go build -o $(APP_NAME)

# Test the application
# test:
# 	go test ./...

# Clean up build artifacts and database file
clean:
	rm -f $(APP_NAME) $(DB_FILE)

# Initialize the database (create if it doesn't exist)
init-db:
	sqlite3 $(DB_FILE) < $(SCHEMA_FILE) || echo "Database already exists"

# Format Go code
fmt:
	go fmt ./...

# Tidy up dependencies
tidy:
	go mod tidy

# Help documentation for Makefile commands
help:
	@echo "Available commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application binary"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Remove build artifacts and database file"
	@echo "  make init-db   - Initialize the SQLite database"
	@echo "  make fmt       - Format Go code"
	@echo "  make tidy      - Tidy up dependencies"
