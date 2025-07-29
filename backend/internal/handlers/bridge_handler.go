package handlers

import (
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

type BridgeHandler struct {
	bridgeService *services.BridgeService
}

func NewBridgeHandler(bridgeService *services.BridgeService) *BridgeHandler {
	return &BridgeHandler{
		bridgeService: bridgeService,
	}
}

// GetBridgeRoutes handles POST /bridge/routes
func (h *BridgeHandler) GetBridgeRoutes(c *fiber.Ctx) error {
	var req services.BridgeRouteRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// Validate request
	if req.FromChain == 0 || req.ToChain == 0 {
		return errors.BadRequest("FromChain and ToChain are required")
	}
	if req.FromToken == "" || req.FromAmount == "" {
		return errors.BadRequest("FromToken and FromAmount are required")
	}
	if req.UserAddress == "" {
		return errors.BadRequest("UserAddress is required")
	}

	// Set default slippage if not provided
	if req.Slippage == 0 {
		req.Slippage = 0.5
	}

	// Get bridge routes
	routes, err := h.bridgeService.GetRoutes(c.Context(), req)
	if err != nil {
		return err
	}

	// Return routes in expected format
	return c.JSON(fiber.Map{
		"routes": routes,
	})
}

// ExecuteBridge handles POST /bridge/execute
func (h *BridgeHandler) ExecuteBridge(c *fiber.Ctx) error {
	var req struct {
		RouteID     string `json:"routeId"`
		UserAddress string `json:"userAddress"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	if req.RouteID == "" || req.UserAddress == "" {
		return errors.BadRequest("RouteID and UserAddress are required")
	}

	// TODO: Implement actual bridge execution
	// For now, return mock transaction hash
	mockTxHash := "0x" + generateMockHash()

	return c.JSON(fiber.Map{
		"txHash": mockTxHash,
	})
}

func generateMockHash() string {
	// Generate a mock 64-character hex string
	const hexChars = "0123456789abcdef"
	hash := make([]byte, 64)
	for i := range hash {
		hash[i] = hexChars[i%16]
	}
	return string(hash)
}