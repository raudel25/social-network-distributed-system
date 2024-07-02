# Go parameters
GOCMD=go
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=socialnetwork

JOIN ?= ""
CONTAINER_PORT ?= 5000

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: deps
deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

.PHONY: dev
dev:
	$(GOCMD) run cmd/main.go -p $(CONTAINER_PORT) -j $(JOIN)

.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto
	protoc --go_out=. --go-grpc_out=. pkg/services/proto/*.proto

# -------------------------------------------- Docker commands -----------------------------------------------------------------------

PORT ?= 5000

.PHONY: docker-build
docker-build:
	docker build -t socialnetwork .

.PHONY: docker-run
docker-run:
	docker run -it --rm -p $(PORT):5000 -v $(PWD):/app socialnetwork

.PHONY: docker-dev
docker-dev:
	docker run -it --rm -p $(PORT):5000 -v $(PWD):/app socialnetwork make dev

.PHONY: docker-proto
docker-proto:
	docker run -it --rm -p $(PORT):5000 -v $(PWD):/app socialnetwork make proto
