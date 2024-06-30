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

// Saves a protobuf message to the Chord network
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

// Loads a protobuf message from the Chord network
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

// Removes a file from the Chord network
func Delete(node *chord.Node, path string) error {
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	err := node.RemoveKey(fullPath)
	if err != nil {
		log.Printf("Error deleting file: %v", err)
		return status.Error(codes.Internal, "Couldn't delete file")
	}
	return nil
}

// Checks if a file exists in the Chord network
func FileExists(node *chord.Node, path string) (bool, error) {
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	log.Printf("Checking if file exists: %s", fullPath)

	_, err := node.GetKey(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, status.Errorf(codes.Internal, "Couldn't get files")
	}

	log.Printf("File already exists: %v", fullPath)
	return true, nil
}
