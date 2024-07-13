GOCMD=go
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Default ID (can be overridden from command line)
ID ?= 0

# Calculate ports based on ID
PORT := $(shell echo $$((10000 + $(ID))))
BL := $(shell echo $$((11000 + $(ID))))
BR := $(shell echo $$((12000 + $(ID))))

# -------------------------------------------- Local commands -------------------------------------------------------------------

# Install dependencies
.PHONY: deps
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

# Run the application in development mode
.PHONY: dev
dev:
	$(GOCMD) run cmd/main.go -p $(PORT) -bl $(BL) -br $(BR)

# Generate chord protocol buffer code
.PHONY: proto-chord
proto-chord:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto

# Generate services protocol buffer code
.PHONY: proto-services
proto-services:
	protoc --go_out=. --go-grpc_out=. internal/services/proto/*.proto

# Generate all protocol buffer code
.PHONY: proto
proto:
	make proto-chord
	make proto-services

# -------------------------------------------- Docker commands -----------------------------------------------------------------------

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t socialnetwork .

# Run Docker container
.PHONY: docker-run
docker-run:
	docker run -it --rm -p $(PORT):10000 -p $(BL):11000 -p $(BR):12000 -v $(PWD):/app socialnetwork

# Run development environment in Docker
.PHONY: docker-dev
docker-dev:
	docker run -it --rm -p $(PORT):10000 -p $(BL):11000 -p $(BR):12000 -v $(PWD):/app socialnetwork make dev

# Generate protocol buffer code in Docker
.PHONY: docker-proto
docker-proto:
	docker run -it --rm -p $(PORT):10000 -p $(BL):11000 -p $(BR):12000 -v $(PWD):/app socialnetwork make proto