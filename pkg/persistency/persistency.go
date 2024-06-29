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
// It takes a Chord node, a protobuf message, and a file path as input
// The file is saved in the "resources" directory with the .bin extension
func Save[T protoreflect.ProtoMessage](node *chord.Node, object T, path string) error {
	// Construct the full file path
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	log.Printf("Saving file: %s", fullPath)

	// Marshal the protobuf message to bytes
	data, err := proto.Marshal(object)
	if err != nil {
		log.Printf("Error serializing object: %v", err)
		return status.Error(codes.Internal, "Error saving data")
	}

	// Encode the bytes to a base64 string
	strData := base64.StdEncoding.EncodeToString(data)

	// Store the base64 string in the Chord network
	err = node.SetKey(fullPath, strData)
	if err != nil {
		log.Println("Error saving file")
		return status.Error(codes.Internal, "Error saving data")
	}

	return nil
}

// Loads a protobuf message from the Chord network
// It takes a Chord node, a file path, and a protobuf message as input
// The file is loaded from the "resources" directory with the .bin extension
func Load[T protoreflect.ProtoMessage](node *chord.Node, path string, result T) (T, error) {
	// Construct the full file path
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	log.Printf("Loading file: %s", fullPath)

	// Create an empty instance of the protobuf message
	var empty T

	// Retrieve the base64 string from the Chord network
	dataStr, err := node.GetKey(fullPath)
	if err != nil {
		log.Printf("Error getting file: %v", err)
		return empty, status.Errorf(codes.Internal, "")
	}

	// Decode the base64 string to bytes
	data, err := base64.StdEncoding.DecodeString(*dataStr)
	if err != nil {
		log.Printf("Error decoding object: %v", err)
		return empty, status.Errorf(codes.Internal, "")
	}

	// Unmarshal the bytes into the protobuf message
	err = proto.Unmarshal(data, result)
	if err != nil {
		log.Printf("Error deserializing object: %v", err)
		return empty, status.Errorf(codes.Internal, "")
	}

	return result, nil
}

// Removes a file from the Chord network
// It takes a Chord node and a file path as input
// The file is removed from the "resources" directory with the .bin extension
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
// It takes a Chord node and a file path as input
// The file is checked in the "resources" directory with the .bin extension
func FileExists(node *chord.Node, path string) bool {
	fullPath := strings.ToLower(filepath.ToSlash(filepath.Join("resources", path+".bin")))

	log.Printf("Checking if file exists: %s", fullPath)

	// Check if the file exists in the Chord network
	if _, err := node.GetKey(fullPath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	log.Printf("File already exists: %v", fullPath)
	return true
}
