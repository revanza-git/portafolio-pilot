package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/internal/middleware"
	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService *services.AuthService
	siweService *services.SIWEService
	jwtSecret   string
	jwtExpiry   int
}

func NewAuthHandler(authService *services.AuthService, siweService *services.SIWEService, jwtSecret string, jwtExpiry int) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		siweService: siweService,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
	}
}

// GetNonce handles POST /auth/siwe/nonce
func (h *AuthHandler) GetNonce(c *fiber.Ctx) error {
	var req NonceRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	if req.Address == "" {
		return errors.BadRequest("Address is required")
	}

	// Validate address format using SIWE service
	if !h.siweService.ValidateAddress(req.Address) {
		return errors.BadRequest("Invalid Ethereum address")
	}

	// Generate nonce
	nonce, err := h.siweService.GenerateNonce(c.Context(), req.Address)
	if err != nil {
		return err
	}

	// Generate SIWE message
	message, err := h.siweService.GenerateSIWEMessage(req.Address, nonce)
	if err != nil {
		return err
	}

	return c.JSON(NonceResponse{
		Nonce:   nonce,
		Message: message,
	})
}

// Verify handles POST /auth/siwe/verify
func (h *AuthHandler) Verify(c *fiber.Ctx) error {
	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// Validate request
	if req.Message == "" || req.Signature == "" {
		return errors.BadRequest("Message and signature are required")
	}

	// Verify SIWE signature
	user, err := h.siweService.VerifySignature(c.Context(), req.Message, req.Signature)
	if err != nil {
		return err
	}

	// Generate JWT token
	token, err := h.generateJWT(user)
	if err != nil {
		return errors.Internal("Failed to generate token")
	}

	return c.JSON(AuthResponse{
		Token:     token,
		ExpiresIn: h.jwtExpiry * 3600, // Convert hours to seconds
		Address:   user.Address,
		User:      user,
	})
}

// GetMe handles GET /auth/me
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	// Get claims from middleware
	claims, ok := c.Locals("claims").(*middleware.Claims)
	if !ok {
		return errors.Unauthorized("Invalid token claims")
	}

	// Get user by address
	user, err := h.authService.GetUserByAddress(c.Context(), claims.Address)
	if err != nil {
		return errors.NotFound("User not found")
	}

	return c.JSON(UserProfileResponse{
		User:    user,
		Address: user.Address,
	})
}

// SendMagicLink handles POST /auth/magic-link (stub implementation)
func (h *AuthHandler) SendMagicLink(c *fiber.Ctx) error {
	var req MagicLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	if req.Email == "" {
		return errors.BadRequest("Email is required")
	}

	// Generate magic link token (stub implementation)
	token := make([]byte, 32)
	rand.Read(token)
	magicToken := hex.EncodeToString(token)

	// TODO: Implement actual email sending
	// For now, just log the magic link
	magicLink := fmt.Sprintf("https://your-domain.com/auth/verify-magic?token=%s", magicToken)
	
	// In development, return the link; in production, just return success
	return c.JSON(MagicLinkResponse{
		Message:   "Magic link sent to your email",
		MagicLink: magicLink, // Remove this in production
	})
}

// generateJWT creates a JWT token for the user
func (h *AuthHandler) generateJWT(user *models.User) (string, error) {
	claims := &middleware.Claims{
		Address: user.Address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(h.jwtExpiry))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// Request/Response types

type NonceRequest struct {
	Address string `json:"address" validate:"required"`
}

type NonceResponse struct {
	Nonce   string `json:"nonce"`
	Message string `json:"message"`
}

type VerifyRequest struct {
	Message   string `json:"message" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int          `json:"expires_in"`
	Address   string       `json:"address"`
	User      *models.User `json:"user"`
}

type UserProfileResponse struct {
	User    *models.User `json:"user"`
	Address string       `json:"address"`
}

type MagicLinkRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type MagicLinkResponse struct {
	Message   string `json:"message"`
	MagicLink string `json:"magic_link,omitempty"` // Only for development
}