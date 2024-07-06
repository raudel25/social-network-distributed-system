# Go parameters
GOCMD=go
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

CONTAINER_PORT ?= 5000
CONTAINER_BL ?= 6000
CONTAINER_BR ?= 7000

.PHONY: deps
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

.PHONY: dev
dev:
	$(GOCMD) run cmd/main.go -p $(CONTAINER_PORT) -bl $(CONTAINER_BL) -br $(CONTAINER_BR)

.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto
	protoc --go_out=. --go-grpc_out=. internal/services/proto/*.proto

# -------------------------------------------- Docker commands -----------------------------------------------------------------------

PORT ?= 5000
BL ?= 6000
BR ?= 7000

.PHONY: docker-build
docker-build:
	docker build -t socialnetwork .

.PHONY: docker-run
docker-run:
	docker run -it --rm -p $(PORT):5000 -p $(BL):6000 -p $(BR):7000 -v $(PWD):/app socialnetwork

.PHONY: docker-dev
docker-dev:
	docker run -it --rm -p $(PORT):5000 -p $(BL):6000 -p $(BR):7000 -v $(PWD):/app socialnetwork make dev

.PHONY: docker-proto
docker-proto:
	docker run -it --rm -p $(PORT):5000 -p $(BL):6000 -p $(BR):7000 -v $(PWD):/app socialnetwork make proto
