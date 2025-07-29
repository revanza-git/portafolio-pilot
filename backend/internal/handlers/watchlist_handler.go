package handlers

import (
	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WatchlistHandler struct {
	watchlistRepo repos.WatchlistRepository
}

func NewWatchlistHandler(watchlistRepo repos.WatchlistRepository) *WatchlistHandler {
	return &WatchlistHandler{
		watchlistRepo: watchlistRepo,
	}
}

// GetWatchlist handles GET /watchlist
func (h *WatchlistHandler) GetWatchlist(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	watchlists, err := h.watchlistRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		logger.Error("Failed to get watchlist",
			"error", err.Error(),
			"userID", userID,
		)
		return errors.Internal("Failed to get watchlist")
	}

	return c.JSON(watchlists)
}

// CreateWatchlistItem handles POST /watchlist
func (h *WatchlistHandler) CreateWatchlistItem(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	var req models.CreateWatchlistRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// Validate item type
	if req.ItemType != models.WatchlistItemTypeToken &&
		req.ItemType != models.WatchlistItemTypePool &&
		req.ItemType != models.WatchlistItemTypeProtocol {
		return errors.BadRequest("Invalid item_type. Must be one of: token, pool, protocol")
	}

	// Validate item_ref_id
	if req.ItemRefID <= 0 {
		return errors.BadRequest("item_ref_id must be a positive integer")
	}

	// Check if item already exists (upsert guard)
	exists, err := h.watchlistRepo.ExistsByUserIDAndItem(c.Context(), userID, req.ItemType, req.ItemRefID)
	if err != nil {
		logger.Error("Failed to check watchlist item existence",
			"error", err.Error(),
			"userID", userID,
			"itemType", req.ItemType,
			"itemRefID", req.ItemRefID,
		)
		return errors.Internal("Failed to check watchlist item")
	}

	if exists {
		return errors.BadRequest("Item already exists in watchlist")
	}

	watchlist := &models.Watchlist{
		UserID:    userID,
		ItemType:  req.ItemType,
		ItemRefID: req.ItemRefID,
	}

	err = h.watchlistRepo.Create(c.Context(), watchlist)
	if err != nil {
		logger.Error("Failed to create watchlist item",
			"error", err.Error(),
			"userID", userID,
			"itemType", req.ItemType,
			"itemRefID", req.ItemRefID,
		)
		return errors.Internal("Failed to create watchlist item")
	}

	return c.Status(201).JSON(watchlist)
}

// DeleteWatchlistItem handles DELETE /watchlist/:id
func (h *WatchlistHandler) DeleteWatchlistItem(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	watchlistIDStr := c.Params("id")
	watchlistID, err := uuid.Parse(watchlistIDStr)
	if err != nil {
		return errors.BadRequest("Invalid watchlist item ID")
	}

	err = h.watchlistRepo.Delete(c.Context(), watchlistID, userID)
	if err != nil {
		logger.Error("Failed to delete watchlist item",
			"error", err.Error(),
			"watchlistID", watchlistID,
			"userID", userID,
		)
		if err.Error() == "watchlist item not found" {
			return errors.NotFound("Watchlist item")
		}
		return errors.Internal("Failed to delete watchlist item")
	}

	return c.SendStatus(204)
}