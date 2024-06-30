package socialnetwork;

import (
	"context"
	"crypto/rsa"
	"errors"
	"log"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/raudel25/social-network-distributed-system/pkg/persistency"
	auth_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_auth"
	users_pb "github.com/raudel25/social-network-distributed-system/pkg/services/grpc_users"
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

func (server *AuthServer) Login(ctx context.Context, request *auth_pb.LoginRequest) (*auth_pb.LoginResponse, error) {
	user, err := loadUser(request.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Wrong username or password")
	}
	if err := verifyPassword(user.PasswordHash, request.Password); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Wrong username or password")
	}
	tokenString, err := server.generateToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}
	return &auth_pb.LoginResponse{Token: tokenString}, nil
}

func (server *AuthServer) SignUp(ctx context.Context, request *auth_pb.SignUpRequest) (*auth_pb.SignUpResponse, error) {
	user := request.GetUser()
	if err := saveUser(user); err != nil {
		return &auth_pb.SignUpResponse{}, err
	}
	return &auth_pb.SignUpResponse{}, nil
}

func ValidateRequest(ctx context.Context) (*jwt.Token, error) {
	publicKey, err := loadPublicKey(rsaPublic)
	if err != nil {
		log.Fatalf("Error reading and parsing the jwt public key: %v", err)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Valid token required.")
	}
	jwtToken, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Valid token required.")
	}
	token, err := validateToken(jwtToken[0], publicKey)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Valid token required.")
	}
	return token, nil
}

func StartAuthServer(network, address string) {
	log.Println("Auth service started")

	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	privateKey, err := loadPrivateKey(rsaPrivate)
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				UnaryLoggingInterceptor,
			),
		), grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				StreamLoggingInterceptor,
			),
		),
	)

	auth_pb.RegisterAuthServer(s, &AuthServer{jwtPrivateKey: privateKey})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loadUser(username string) (*users_pb.User, error) {
	user := &users_pb.User{}
	path := filepath.Join("User", strings.ToLower(username))
	return persistency.Load(node, path, user)

}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func saveUser(user *users_pb.User) error {
	user.Username = strings.ToLower(user.Username)
	path := filepath.Join("User", user.Username)
	if persistency.FileExists(node, path) {
		return status.Error(codes.AlreadyExists, "Username is taken")
	}
	return persistency.Save(node, user, path)
}

func (server *AuthServer) generateToken(user *users_pb.User) (string, error) {
	claims := jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"iss":   "auth.service",
		"iat":   time.Now().Unix(),
		"email": user.Email,
		"sub":   user.Username,
		"name":  user.Name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(server.jwtPrivateKey)
}

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token")
		}
		return publicKey, nil
	})
}

func extractUsernameFromToken(token *jwt.Token) (string, error) {
	username, ok := token.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return "", status.Errorf(codes.Internal, "Failed to extract username from token")
	}
	return username, nil
}

func checkPermission(ctx context.Context, requestedUsername string) error {
	token, err := ValidateRequest(ctx)
	if err != nil {
		return err
	}

	username, err := extractUsernameFromToken(token)
	if err != nil {
		return err
	}

	if username != requestedUsername {
		return status.Errorf(codes.PermissionDenied, "You are not authorized to edit this user")
	}
	return nil
}
