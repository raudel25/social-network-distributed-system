package socialnetwork

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}
