package services

import (
	"context"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	auth_pb "github.com/raudel25/social-network-distributed-system/pkg/services/auth"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/users"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	*auth_pb.UnimplementedAuthServer
	jwtPrivateKey *rsa.PrivateKey
}

func (server *AuthServer) Login(_ context.Context, request *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error) {

	user := &users_pb.User{}
	user, err := persistency.Load(node, filepath.Join("User", request.GetUsername()), user)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Wrong username or password")
	}

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["iss"] = "auth.service"
	claims["iat"] = time.Now().Unix()
	claims["email"] = user.Email
	claims["sub"] = user.Username
	claims["name"] = user.Name

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(server.jwtPrivateKey)
	if err != nil {
		// fmt.Println("Error creating token: %v", err)
		return nil, status.Errorf(codes.Internal, "")
	}

	return &auth_pb.LoginResponse{Token: tokenString}, nil
}

func (*AuthServer) SignUp(_ context.Context, request *auth_pb.SignUpRequest) (*auth_pb.SignUpResponse, error) {
	user := request.GetUser()
	user.Username = strings.ToLower(user.Username)
	path := filepath.Join("User", user.Username)

	if persistency.FileExists(node, path) {
		return &auth_pb.SignUpResponse{}, status.Error(codes.AlreadyExists, "Username is taken")
	}

	err := persistency.Save(node, user, path)

	if err != nil {
		return &auth_pb.SignUpResponse{}, err
	}

	return &auth_pb.SignUpResponse{}, nil
}

func StartAuthServer(network string, address string) {
	log.Println("Auth service started")

	lis, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	key, err := ioutil.ReadFile(rsaPrivate)
	if err != nil {
		log.Fatalf("Error reading the jwt private key: %v", err)
	}
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		log.Fatalf("Error parsing the jwt private key: %v", err)
	}

	s := grpc.NewServer()

	auth_pb.RegisterAuthServer(s, &AuthServer{jwtPrivateKey: parsedKey})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			// log.Println("Unexpected signing method: %v", t.Header["alg"])
			return nil, errors.New("invalid token")
		}
		return publicKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}

func ValidateRequest(ctx context.Context) (*jwt.Token, error) {
	var (
		token *jwt.Token
		err   error
	)

	key, err := ioutil.ReadFile(rsaPublic)

	if err != nil {
		log.Fatalf("Error reading the jwt public key: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(key)

	if err != nil {
		log.Fatalf("Error parsing the jwt public key: %v", err)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "valid token required.")
	}

	jwtToken, ok := md["authorization"]

	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "valid token required.")
	}

	token, err = validateToken(jwtToken[0], publicKey)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "valid token required.")
	}

	return token, nil
}
