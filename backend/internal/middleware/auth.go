package middleware

import (
	"context"
	"strings"

	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.Unauthorized("Missing authorization header")
		}

		// Check Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return errors.Unauthorized("Invalid authorization header format")
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Unauthorized("Invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil {
			return errors.Unauthorized("Invalid token")
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			return errors.Unauthorized("Invalid token claims")
		}

		// Store user info in context
		c.Locals("address", claims.Address)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// OptionalAuth is like JWTAuth but doesn't fail if no token is provided
func OptionalAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := tokenParts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*Claims); ok {
				c.Locals("address", claims.Address)
				c.Locals("claims", claims)
			}
		}

		return c.Next()
	}
}

// JWTAuthWithUser extends JWTAuth to also resolve user information
func JWTAuthWithUser(secret string, userRepo repos.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.Unauthorized("Missing authorization header")
		}

		// Check Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return errors.Unauthorized("Invalid authorization header format")
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Unauthorized("Invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil {
			return errors.Unauthorized("Invalid token")
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			return errors.Unauthorized("Invalid token claims")
		}

		// Get user by address to resolve userID and admin status
		user, err := userRepo.GetByAddress(context.Background(), claims.Address)
		if err != nil {
			return errors.Unauthorized("User not found")
		}

		// Store user info in context
		c.Locals("address", claims.Address)
		c.Locals("userID", user.ID)
		c.Locals("isAdmin", user.IsAdmin)
		c.Locals("claims", claims)
		c.Locals("user", user)

		return c.Next()
	}
}

// AdminAuth middleware checks if user is admin
func AdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("isAdmin").(bool)
		if !ok || !isAdmin {
			return errors.Forbidden("Admin access required")
		}
		return c.Next()
	}
}