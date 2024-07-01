.PHONY: dev
dev:
	go run cmd/main.go

.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. pkg/chord/grpc/chord.proto
	protoc --go_out=. --go-grpc_out=. pkg/services/proto/auth.proto
	protoc --go_out=. --go-grpc_out=. pkg/services/proto/db_models.proto
	protoc --go_out=. --go-grpc_out=. pkg/services/proto/follow.proto
	protoc --go_out=. --go-grpc_out=. pkg/services/proto/posts.proto
	protoc --go_out=. --go-grpc_out=. pkg/services/proto/users.proto