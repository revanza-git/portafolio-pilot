package handlers

import (
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

type SwapHandler struct {
	swapService *services.SwapService
}

func NewSwapHandler(swapService *services.SwapService) *SwapHandler {
	return &SwapHandler{
		swapService: swapService,
	}
}

// GetSwapQuote handles POST /swap/quote
func (h *SwapHandler) GetSwapQuote(c *fiber.Ctx) error {
	var req services.SwapQuoteRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// Validate request
	if req.ChainID == 0 {
		return errors.BadRequest("ChainID is required")
	}
	if req.FromToken == "" || req.ToToken == "" {
		return errors.BadRequest("FromToken and ToToken are required")
	}
	if req.FromAmount == "" {
		return errors.BadRequest("FromAmount is required")
	}
	if req.UserAddress == "" {
		return errors.BadRequest("UserAddress is required")
	}

	// Set default slippage if not provided
	if req.Slippage == 0 {
		req.Slippage = 0.5
	}

	// Get swap quotes
	quotes, err := h.swapService.GetQuotes(c.Context(), req)
	if err != nil {
		return err
	}

	// Return quotes as array (frontend expects array directly)
	return c.JSON(quotes)
}

// ExecuteSwap handles POST /swap/execute
func (h *SwapHandler) ExecuteSwap(c *fiber.Ctx) error {
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

	// TODO: Implement actual swap execution
	// For now, return mock transaction hash
	mockTxHash := "0x" + generateSwapTxHash()

	return c.JSON(fiber.Map{
		"txHash": mockTxHash,
	})
}

func generateSwapTxHash() string {
	// Generate a mock 64-character hex string
	const hexChars = "0123456789abcdef"
	hash := make([]byte, 64)
	for i := range hash {
		hash[i] = hexChars[(i*7)%16] // Different pattern than bridge
	}
	return string(hash)
}