p ?= 5000

.PHONY: dev
dev:
	go run cmd/main.go -p $(p)

.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto