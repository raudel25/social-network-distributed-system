.PHONY: dev
dev:
	go run cmd/main.go

.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto
	protoc --go_out=. --go-grpc_out=. pkg/services/proto/*.proto