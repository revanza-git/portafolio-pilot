package handlers

import (
	"strconv"

	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type YieldHandler struct {
	yieldService *services.YieldService
}

func NewYieldHandler(yieldService *services.YieldService) *YieldHandler {
	return &YieldHandler{
		yieldService: yieldService,
	}
}

// GetYieldPools handles GET /yield/pools
func (h *YieldHandler) GetYieldPools(c *fiber.Ctx) error {
	// Parse query parameters
	filters := repos.YieldPoolFilters{
		Chain:        getStringParam(c, "chain"),
		ChainID:      getIntParam(c, "chainId"),
		MinTVL:       getFloat64Param(c, "minTvl"),
		MinAPY:       getFloat64Param(c, "minApy"),
		ProtocolSlug: getStringParam(c, "protocol"),
		RiskLevel:    getStringParam(c, "riskLevel"),
		IsActive:     getBoolParam(c, "active"),
		SortBy:       c.Query("sort", "apy"),
		Limit:        getIntValueOrDefault(c, "limit", 20),
		Offset:       getIntValueOrDefault(c, "offset", 0),
	}

	// Validate sort parameter
	validSorts := map[string]bool{
		"apy": true, "tvl": true, "name": true,
	}
	if !validSorts[filters.SortBy] {
		filters.SortBy = "apy"
	}

	// Get pools from service
	pools, total, err := h.yieldService.GetPools(c.Context(), filters)
	if err != nil {
		return err
	}

	// Calculate pagination metadata
	totalPages := (int(total) + filters.Limit - 1) / filters.Limit
	page := (filters.Offset / filters.Limit) + 1

	return c.JSON(fiber.Map{
		"data": pools,
		"meta": fiber.Map{
			"page":       page,
			"limit":      filters.Limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetYieldPositions handles GET /yield/positions/:address
func (h *YieldHandler) GetYieldPositions(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return errors.BadRequest("Address parameter is required")
	}

	// Validate Ethereum address format
	if !isValidEthereumAddress(address) {
		return errors.BadRequest("Invalid Ethereum address format")
	}

	// Parse query parameters
	filters := repos.PositionFilters{
		IsActive: getBoolParam(c, "active"),
		ChainID:  getIntParam(c, "chainId"),
		Limit:    getIntValueOrDefault(c, "limit", 50),
		Offset:   getIntValueOrDefault(c, "offset", 0),
	}

	// Get positions from service
	summary, err := h.yieldService.GetUserPositions(c.Context(), address, filters)
	if err != nil {
		return err
	}

	return c.JSON(summary)
}

// ClaimRewards handles POST /yield/positions/:address/:positionId/claim
func (h *YieldHandler) ClaimRewards(c *fiber.Ctx) error {
	address := c.Params("address")
	positionIDStr := c.Params("positionId")

	if address == "" || positionIDStr == "" {
		return errors.BadRequest("Address and position ID parameters are required")
	}

	// Validate Ethereum address format
	if !isValidEthereumAddress(address) {
		return errors.BadRequest("Invalid Ethereum address format")
	}

	// Parse position UUID
	positionID, err := uuid.Parse(positionIDStr)
	if err != nil {
		return errors.BadRequest("Invalid position ID format")
	}

	// Claim rewards through service
	response, err := h.yieldService.ClaimRewards(c.Context(), address, positionID)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// GetYieldPoolsByProtocol handles GET /yield/pools/protocol/:slug
func (h *YieldHandler) GetYieldPoolsByProtocol(c *fiber.Ctx) error {
	protocolSlug := c.Params("slug")
	if protocolSlug == "" {
		return errors.BadRequest("Protocol slug parameter is required")
	}

	// Get protocol first
	protocol, err := h.yieldService.GetProtocolBySlug(c.Context(), protocolSlug)
	if err != nil {
		return err
	}

	// Get pools for this protocol
	activeOnly := getBoolValueOrDefault(c, "active", true)
	pools, err := h.yieldService.GetPoolsByProtocol(c.Context(), protocol.ID, activeOnly)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"protocol": protocol,
		"pools":    pools,
	})
}

// GetYieldPoolsByChain handles GET /yield/pools/chain/:chainId
func (h *YieldHandler) GetYieldPoolsByChain(c *fiber.Ctx) error {
	chainIDStr := c.Params("chainId")
	if chainIDStr == "" {
		return errors.BadRequest("Chain ID parameter is required")
	}

	chainID, err := strconv.Atoi(chainIDStr)
	if err != nil {
		return errors.BadRequest("Invalid chain ID format")
	}

	// Get pools for this chain
	pools, err := h.yieldService.GetPoolsByChain(c.Context(), chainID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"chainId": chainID,
		"pools":   pools,
	})
}

// GetTopYieldPools handles GET /yield/pools/top
func (h *YieldHandler) GetTopYieldPools(c *fiber.Ctx) error {
	limit := getIntValueOrDefault(c, "limit", 10)
	if limit > 100 {
		limit = 100 // Cap at 100
	}

	pools, err := h.yieldService.GetTopPoolsByTVL(c.Context(), limit)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"pools": pools,
	})
}

// GetProtocols handles GET /yield/protocols
func (h *YieldHandler) GetProtocols(c *fiber.Ctx) error {
	// Parse query parameters
	filters := repos.ProtocolFilters{
		Category:  getStringParam(c, "category"),
		IsActive:  getBoolParam(c, "active"),
		RiskLevel: getStringParam(c, "riskLevel"),
		SortBy:    c.Query("sort", "tvl"),
		Limit:     getIntValueOrDefault(c, "limit", 20),
		Offset:    getIntValueOrDefault(c, "offset", 0),
	}

	// Validate sort parameter
	validSorts := map[string]bool{
		"name": true, "tvl": true, "category": true,
	}
	if !validSorts[filters.SortBy] {
		filters.SortBy = "tvl"
	}

	// Get protocols from service
	protocols, total, err := h.yieldService.GetProtocols(c.Context(), filters)
	if err != nil {
		return err
	}

	// Calculate pagination metadata
	totalPages := (int(total) + filters.Limit - 1) / filters.Limit
	page := (filters.Offset / filters.Limit) + 1

	return c.JSON(fiber.Map{
		"data": protocols,
		"meta": fiber.Map{
			"page":       page,
			"limit":      filters.Limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// CreatePosition handles POST /yield/positions/:address (internal/admin use)
func (h *YieldHandler) CreatePosition(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return errors.BadRequest("Address parameter is required")
	}

	// Validate Ethereum address format
	if !isValidEthereumAddress(address) {
		return errors.BadRequest("Invalid Ethereum address format")
	}

	var req services.CreatePositionRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// Create position through service
	position, err := h.yieldService.CreatePosition(c.Context(), address, req)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(position)
}

// UpdatePosition handles PUT /yield/positions/:positionId (internal/admin use)
func (h *YieldHandler) UpdatePosition(c *fiber.Ctx) error {
	positionIDStr := c.Params("positionId")
	if positionIDStr == "" {
		return errors.BadRequest("Position ID parameter is required")
	}

	// Parse position UUID
	positionID, err := uuid.Parse(positionIDStr)
	if err != nil {
		return errors.BadRequest("Invalid position ID format")
	}

	var req services.UpdatePositionRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// Update position through service
	position, err := h.yieldService.UpdatePosition(c.Context(), positionID, req)
	if err != nil {
		return err
	}

	return c.JSON(position)
}

// Helper functions

func getStringParam(c *fiber.Ctx, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	return &value
}

func getIntParam(c *fiber.Ctx, key string) *int {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return &intValue
	}
	return nil
}

func getFloat64Param(c *fiber.Ctx, key string) *float64 {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return &floatValue
	}
	return nil
}

func getBoolParam(c *fiber.Ctx, key string) *bool {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	if boolValue, err := strconv.ParseBool(value); err == nil {
		return &boolValue
	}
	return nil
}

func getIntValueOrDefault(c *fiber.Ctx, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

func getBoolValueOrDefault(c *fiber.Ctx, key string, defaultValue bool) bool {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue
	}
	return defaultValue
}

func isValidEthereumAddress(address string) bool {
	// Basic validation for Ethereum address format
	if len(address) != 42 {
		return false
	}
	if address[:2] != "0x" {
		return false
	}
	// Additional validation could be added here
	return true
}