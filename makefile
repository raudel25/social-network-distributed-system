# Go parameters
GOCMD=go
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

CONTAINER_BROAD ?= 8000
CONTAINER_PORT ?= 5000

.PHONY: deps
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

.PHONY: dev
dev:
	$(GOCMD) run cmd/main.go -p $(CONTAINER_PORT) -b $(CONTAINER_BROAD)

.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto
	protoc --go_out=. --go-grpc_out=. internal/services/proto/*.proto

# -------------------------------------------- Docker commands -----------------------------------------------------------------------

PORT ?= 5000
BROAD ?= 8000

.PHONY: docker-build
docker-build:
	docker build -t socialnetwork .

.PHONY: docker-run
docker-run:
	docker run -it --rm -p $(PORT):5000 -p $(BROAD):8000 -v $(PWD):/app socialnetwork

.PHONY: docker-dev
docker-dev:
	docker run -it --rm -p $(PORT):5000 -p $(BROAD):8000 -v $(PWD):/app socialnetwork make dev

.PHONY: docker-proto
docker-proto:
	docker run -it --rm -p $(PORT):5000 -p $(BROAD):8000 -v $(PWD):/app socialnetwork make proto
