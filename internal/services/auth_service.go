package socialnetwork

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	socialnetwork_pb "github.com/raudel25/social-network-distributed-system/internal/services/grpc"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	*socialnetwork_pb.UnimplementedAuthServer
	jwtPrivateKey *rsa.PrivateKey
}

func (server *AuthServer) Login(ctx context.Context, request *socialnetwork_pb.LoginRequest) (*socialnetwork_pb.LoginResponse, error) {
	username := request.GetUsername()

	user, err := loadUser(username)

	if err != nil {
		return nil, err
	}

	if err := verifyPassword(user.PasswordHash, request.Password); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Wrong username or password")
	}

	tokenString, err := server.generateToken(user)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &socialnetwork_pb.LoginResponse{Token: tokenString}, nil
}

func (server *AuthServer) SignUp(ctx context.Context, request *socialnetwork_pb.SignUpRequest) (*socialnetwork_pb.SignUpResponse, error) {
	user := request.GetUser()

	if !isEmailValid(user.Email) {
		return &socialnetwork_pb.SignUpResponse{}, status.Errorf(codes.InvalidArgument, "Invalid email")
	}

	if exists, err := existsUser(user.Username); exists || err != nil {
		return &socialnetwork_pb.SignUpResponse{}, status.Errorf(codes.InvalidArgument, "Fail to sign up")
	}

	if err := saveUser(user); err != nil {
		return &socialnetwork_pb.SignUpResponse{}, err
	}

	return &socialnetwork_pb.SignUpResponse{}, nil
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

	socialnetwork_pb.RegisterAuthServer(s, &AuthServer{jwtPrivateKey: privateKey})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// ===============================================================================================================================================

func validateRequest(ctx context.Context) (*jwt.Token, error) {
	publicKey, err := loadPublicKey(rsaPublic)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error reading and parsing the jwt public key: %v", err)
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
		log.Printf("Error validating token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "Valid token required.")
	}
	return token, nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (server *AuthServer) generateToken(user *socialnetwork_pb.User) (string, error) {
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
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("malformed token: %v", err)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, fmt.Errorf("token is either expired or not active yet: %v", err)
			} else {
				return nil, fmt.Errorf("couldn't handle this token: %v", err)
			}
		}
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return parsedToken, nil
}

func extractUsernameFromToken(token *jwt.Token) (string, error) {
	username, ok := token.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return "", status.Errorf(codes.Internal, "Failed to extract username from token")
	}
	return username, nil
}

func checkPermission(ctx context.Context, requestedUsername string) error {
	token, err := validateRequest(ctx)
	if err != nil {
		return err
	}

	username, err := extractUsernameFromToken(token)
	if err != nil {
		return err
	}

	if username != requestedUsername {
		return status.Errorf(codes.PermissionDenied, "You are not authorized to perform this action")
	}
	return nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	return parsedKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
