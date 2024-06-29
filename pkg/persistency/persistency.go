package persistency

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"log"

	"github.com/raudel25/social-network-distributed-system/pkg/chord"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Save[T protoreflect.ProtoMessage](node *chord.Node, object T, path string) error {
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	log.Printf("Saving file: %s", fullPath)

	data, err := proto.Marshal(object)

	if err != nil {
		log.Printf("Error serializing object: %v", err)
		return status.Error(codes.Internal, "Error saving data")
	}
	
	strData := base64.StdEncoding.EncodeToString(data)
	err = node.SetKey(fullPath, strData)

	if err != nil {
		log.Println("Error saving file")
		return status.Error(codes.Internal, "Error saving data")
	}

	return nil
}

func Load[T protoreflect.ProtoMessage](node *chord.Node, path string, result T) (T, error) {

	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	log.Printf("Loading file: %s", fullPath)

	var empty T

	dataStr, err := node.GetKey(fullPath)

	if err != nil {
		log.Printf("Error getting file: %v", err)
		return empty, status.Errorf(codes.Internal, "")
	}

	data, err := base64.StdEncoding.DecodeString(*dataStr)

	if err != nil {
		log.Printf("Error decoding object: %v", err)
		return empty, status.Errorf(codes.Internal, "")
	}

	err = proto.Unmarshal(data, result)

	if err != nil {
		log.Printf("Error deserializing object: %v", err)
		return empty, status.Errorf(codes.Internal, "")
	}

	return result, nil
}

func Delete(node *chord.Node, path string) error {
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	err := node.RemoveKey(fullPath)

	if err != nil {
		log.Printf("Error deleting file: %v", err)
		return status.Error(codes.Internal, "Couldn't delete file")
	}
	return nil
}

func FileExists(node *chord.Node, path string) bool {
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	log.Printf("Checking if file exists: %s", fullPath)

	if _, err := node.GetKey(fullPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	log.Printf("File already exists: %v", fullPath)
	return true
}
