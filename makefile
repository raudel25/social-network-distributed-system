# Go parameters
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

.PHONY: deps
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

.PHONY: dev
dev:
	$(GOCMD) run cmd/main.go -p $(PORT) -bl $(BL) -br $(BR)

.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto
	protoc --go_out=. --go-grpc_out=. internal/services/proto/*.proto

# -------------------------------------------- Docker commands -------------------------------------------------------------------

.PHONY: docker-build
docker-build:
	docker build -t socialnetwork .

.PHONY: docker-run
docker-run:
	docker run -it --rm -p $(PORT):10000 -p $(BL):11000 -p $(BR):12000 -v $(PWD):/app socialnetwork

.PHONY: docker-dev
docker-dev:
	docker run -it --rm -p $(PORT):10000 -p $(BL):11000 -p $(BR):12000 -v $(PWD):/app socialnetwork make dev

.PHONY: docker-proto
docker-proto:
	docker run -it --rm -p $(PORT):10000 -p $(BL):11000 -p $(BR):12000 -v $(PWD):/app socialnetwork make proto