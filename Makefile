.PHONY: run build test docs clean

# Default binary output
BINARY_NAME?=services-api

# Default directory
BUILD_DIR?=./bin

# Main Go package
MAIN_PACKAGE=.

# Default environment for development
export PORT?=8080
export ENVIRONMENT?=development
export SERVER_HOST?=localhost:8080

# Default targets
all: clean build

# Run the application
run:
	go run ${MAIN_PACKAGE}

# Build the application
build:
	mkdir -p ${BUILD_DIR}
	go build -tags "${BUILD_TAGS}" -ldflags "${LD_FLAGS}" -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PACKAGE}

# Run tests
test:
	go test -v ./...

# Run lint checks
lint:
	go vet ./...
	go fmt ./...

# Generate Swagger documentation
docs:
	swag init -g main.go -o ./docs

# Clean build artifacts
clean:
	rm -rf ${BUILD_DIR}

# Setup development environment
setup:
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

# Build and run with Docker
docker-build:
	docker build -t ${BINARY_NAME}:latest .

docker-run:
	docker run -p ${PORT}:${PORT} -e PORT=${PORT} -e ENVIRONMENT=${ENVIRONMENT} ${BINARY_NAME}:latest

# Run with Docker Compose
docker-compose-up:
	docker-compose up

docker-compose-down:
	docker-compose down

# Load sample data into Postgres via docker-compose
db-load:
	docker-compose exec -T db psql -U kong -d kong_db -f ./scripts/sample_data.sql
