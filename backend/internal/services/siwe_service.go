package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spruceid/siwe-go"
)

type SIWEService struct {
	userRepo  repos.UserRepository
	nonceRepo repos.NonceRepository
	domain    string
}

func NewSIWEService(userRepo repos.UserRepository, nonceRepo repos.NonceRepository, domain string) *SIWEService {
	return &SIWEService{
		userRepo:  userRepo,
		nonceRepo: nonceRepo,
		domain:    domain,
	}
}

// GenerateNonce creates a new nonce for SIWE authentication
func (s *SIWEService) GenerateNonce(ctx context.Context, address string) (string, error) {
	// Validate Ethereum address
	if !common.IsHexAddress(address) {
		return "", errors.BadRequest("Invalid Ethereum address")
	}

	// Normalize address to checksum format
	address = common.HexToAddress(address).Hex()

	// Generate random nonce
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Internal("Failed to generate nonce")
	}
	nonce := hex.EncodeToString(b)

	// Store nonce with 10 minute expiration
	expiresAt := time.Now().Add(10 * time.Minute)
	if err := s.nonceRepo.Store(ctx, address, nonce, expiresAt); err != nil {
		return "", errors.Internal("Failed to store nonce")
	}

	return nonce, nil
}

// VerifySignature verifies a SIWE signature and returns user info
func (s *SIWEService) VerifySignature(ctx context.Context, message, signature string) (*models.User, error) {
	// Parse SIWE message
	siweMessage, err := siwe.ParseMessage(message)
	if err != nil {
		return nil, errors.BadRequest("Invalid SIWE message format")
	}

	// Validate domain
	messageDomain := siweMessage.GetDomain()
	if messageDomain != s.domain {
		// Log the domain mismatch for debugging
		fmt.Printf("Domain mismatch: expected '%s', got '%s'\n", s.domain, messageDomain)
		return nil, errors.BadRequest(fmt.Sprintf("Invalid domain in SIWE message: expected '%s', got '%s'", s.domain, messageDomain))
	}

	// Validate address format
	address := siweMessage.GetAddress().Hex()
	if !common.IsHexAddress(address) {
		return nil, errors.BadRequest("Invalid address in SIWE message")
	}

	// Verify nonce exists and is valid
	nonce := siweMessage.GetNonce()
	isValid, err := s.nonceRepo.ValidateAndUse(ctx, address, nonce)
	if err != nil {
		return nil, errors.Internal("Failed to validate nonce")
	}
	if !isValid {
		return nil, errors.Unauthorized("Invalid or expired nonce")
	}

	// Verify signature
	if err := s.verifyEthereumSignature(message, signature, address); err != nil {
		return nil, errors.Unauthorized("Invalid signature")
	}

	// Get or create user
	user, err := s.userRepo.GetByAddress(ctx, address)
	if err != nil {
		// User not found, create new one
		user, err = s.userRepo.Create(ctx, address, "")
		if err != nil {
			return nil, errors.Internal("Failed to create user")
		}
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID, now); err != nil {
		// Log error but don't fail auth
		// logger.Error("Failed to update last login", err)
	}

	return user, nil
}

// verifyEthereumSignature verifies an Ethereum signature
func (s *SIWEService) verifyEthereumSignature(message, signature, expectedAddress string) error {
	// Remove 0x prefix from signature if present
	if strings.HasPrefix(signature, "0x") {
		signature = signature[2:]
	}

	// Decode signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signature format: %w", err)
	}

	if len(sigBytes) != 65 {
		return fmt.Errorf("signature must be 65 bytes long")
	}

	// Ethereum signatures have recovery ID as last byte, but go-ethereum expects it in different format
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Hash the message with Ethereum prefix
	messageHash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// Recover public key from signature
	pubKey, err := crypto.SigToPub(messageHash.Bytes(), sigBytes)
	if err != nil {
		return fmt.Errorf("failed to recover public key: %w", err)
	}

	// Get address from public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)

	// Compare addresses (case-insensitive)
	if !strings.EqualFold(recoveredAddress.Hex(), expectedAddress) {
		return fmt.Errorf("signature address mismatch: expected %s, got %s", expectedAddress, recoveredAddress.Hex())
	}

	return nil
}

// GenerateSIWEMessage creates a SIWE message for signing
func (s *SIWEService) GenerateSIWEMessage(address, nonce string) (string, error) {
	// Validate address
	if !common.IsHexAddress(address) {
		return "", errors.BadRequest("Invalid Ethereum address")
	}

	// Create SIWE message using the correct constructor
	message, err := siwe.InitMessage(
		s.domain,
		address,
		fmt.Sprintf("https://%s", s.domain),
		nonce,
		map[string]interface{}{
			"statement": "Sign in to DeFi Portfolio Dashboard",
			"version":   "1",
			"chainId":   1,
			"issuedAt":  time.Now().Format(time.RFC3339),
		},
	)
	if err != nil {
		return "", errors.BadRequest("Failed to create SIWE message")
	}

	return message.String(), nil
}

// ValidateAddress validates an Ethereum address format
func (s *SIWEService) ValidateAddress(address string) bool {
	return common.IsHexAddress(address)
}