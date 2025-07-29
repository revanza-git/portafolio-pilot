package handlers

import (
	"strconv"

	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

type PortfolioHandler struct {
	portfolioService *services.PortfolioService
}

func NewPortfolioHandler(portfolioService *services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
	}
}

// GetBalances handles GET /portfolio/:address/balances
func (h *PortfolioHandler) GetBalances(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return errors.BadRequest("Address is required")
	}

	// Parse query parameters
	var chainID *int
	if chainParam := c.Query("chainId"); chainParam != "" {
		chain, err := strconv.Atoi(chainParam)
		if err != nil {
			return errors.BadRequest("Invalid chainId")
		}
		chainID = &chain
	}

	hideSmall := c.Query("hideSmall") == "true"

	// Get balances
	balances, err := h.portfolioService.GetBalances(c.Context(), address, chainID, hideSmall)
	if err != nil {
		return err
	}

	return c.JSON(balances)
}

// GetHistory handles GET /portfolio/:address/history
func (h *PortfolioHandler) GetHistory(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return errors.BadRequest("Address is required")
	}

	// Parse query parameters
	var chainID *int
	if chainParam := c.Query("chainId"); chainParam != "" {
		chain, err := strconv.Atoi(chainParam)
		if err != nil {
			return errors.BadRequest("Invalid chainId")
		}
		chainID = &chain
	}

	period := c.Query("period", "1w")
	interval := c.Query("interval", "1d")

	// Validate period
	validPeriods := map[string]bool{
		"1d": true, "1w": true, "1m": true, 
		"3m": true, "1y": true, "all": true,
	}
	if !validPeriods[period] {
		return errors.BadRequest("Invalid period")
	}

	// Validate interval
	validIntervals := map[string]bool{
		"1h": true, "1d": true, "1w": true,
	}
	if !validIntervals[interval] {
		return errors.BadRequest("Invalid interval")
	}

	// Get history
	history, err := h.portfolioService.GetHistory(c.Context(), address, chainID, period, interval)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"history": history,
	})
}