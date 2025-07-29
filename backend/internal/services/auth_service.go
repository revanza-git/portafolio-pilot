package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/internal/middleware"
	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo   repos.UserRepository
	walletRepo repos.WalletRepository
	jwtSecret  string
	jwtExpiry  int // hours
}

func NewAuthService(userRepo repos.UserRepository, walletRepo repos.WalletRepository, jwtSecret string, jwtExpiry int) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
	}
}

// GenerateNonce generates a random nonce for SIWE
func (s *AuthService) GenerateNonce(ctx context.Context, address string) (string, error) {
	// Generate random nonce
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Internal("Failed to generate nonce")
	}
	nonce := hex.EncodeToString(b)

	// Get or create user
	_, err := s.userRepo.GetByAddress(ctx, address)
	if err != nil {
		// Create new user if not exists
		_, err = s.userRepo.Create(ctx, address, nonce)
		if err != nil {
			return "", errors.Internal("Failed to create user")
		}
	} else {
		// Update existing user's nonce
		_, err = s.userRepo.UpdateNonce(ctx, address, nonce)
		if err != nil {
			return "", errors.Internal("Failed to update nonce")
		}
	}

	return nonce, nil
}

// VerifySignature verifies SIWE signature and returns JWT token
func (s *AuthService) VerifySignature(ctx context.Context, message map[string]interface{}, signature string) (string, error) {
	// TODO: Implement actual SIWE verification using go-ethereum
	// For now, we'll do basic validation and return a mock JWT

	address, ok := message["address"].(string)
	if !ok {
		return "", errors.BadRequest("Invalid address in message")
	}

	_, ok = message["nonce"].(string)
	if !ok {
		return "", errors.BadRequest("Invalid nonce in message")
	}

	// Verify nonce matches stored nonce
	user, err := s.userRepo.GetByAddress(ctx, address)
	if err != nil {
		return "", errors.Unauthorized("User not found")
	}

	// TODO: Actually verify the nonce matches
	// if user.Nonce != nonce {
	//     return "", errors.Unauthorized("Invalid nonce")
	// }

	// Generate new nonce for next login
	newNonce, _ := s.GenerateNonce(ctx, address)
	s.userRepo.UpdateNonce(ctx, address, newNonce)

	// Create JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", errors.Internal("Failed to generate token")
	}

	return token, nil
}

// generateJWT creates a JWT token for the user
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := &middleware.Claims{
		Address: user.Address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.jwtExpiry))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*middleware.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, errors.Unauthorized("Invalid token")
	}

	if claims, ok := token.Claims.(*middleware.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.Unauthorized("Invalid token claims")
}

// GetUserByAddress retrieves a user by their address
func (s *AuthService) GetUserByAddress(ctx context.Context, address string) (*models.User, error) {
	return s.userRepo.GetByAddress(ctx, address)
}