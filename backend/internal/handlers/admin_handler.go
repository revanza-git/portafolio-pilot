package handlers

import (
	"strconv"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userRepo          repos.UserRepository
	featureFlagRepo   repos.FeatureFlagRepository
	systemBannerRepo  repos.SystemBannerRepository
}

func NewAdminHandler(userRepo repos.UserRepository, featureFlagRepo repos.FeatureFlagRepository, systemBannerRepo repos.SystemBannerRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:         userRepo,
		featureFlagRepo:  featureFlagRepo,
		systemBannerRepo: systemBannerRepo,
	}
}

// GetUsers handles GET /admin/users (paginated)
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	// Parse pagination parameters
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	// TODO: Implement actual user listing with pagination
	// For now, return empty response with proper structure
	users := []models.User{}

	return c.JSON(fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"page":       (offset / limit) + 1,
			"limit":      limit,
			"total":      len(users),
			"totalPages": (len(users) + limit - 1) / limit,
		},
	})
}

// GetErrors handles GET /admin/errors (if logged)
func (h *AdminHandler) GetErrors(c *fiber.Ctx) error {
	// TODO: Implement error log retrieval
	// This would depend on how errors are logged in the system
	errors := []interface{}{}

	return c.JSON(fiber.Map{
		"data": errors,
		"meta": fiber.Map{
			"total": len(errors),
		},
	})
}

// GetFeatureFlags handles GET /admin/feature-flags
func (h *AdminHandler) GetFeatureFlags(c *fiber.Ctx) error {
	flags, err := h.featureFlagRepo.GetAll(c.Context())
	if err != nil {
		logger.Error("Failed to get feature flags",
			"error", err.Error(),
		)
		return errors.Internal("Failed to get feature flags")
	}

	return c.JSON(flags)
}

// CreateFeatureFlag handles POST /admin/feature-flags
func (h *AdminHandler) CreateFeatureFlag(c *fiber.Ctx) error {
	var req models.CreateFeatureFlagRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// TODO: Add validation using validator package
	if req.Name == "" {
		return errors.BadRequest("Feature flag name is required")
	}

	if req.Value == nil {
		return errors.BadRequest("Feature flag value is required")
	}

	flag := &models.FeatureFlag{
		Name:  req.Name,
		Value: req.Value,
	}

	err := h.featureFlagRepo.Upsert(c.Context(), flag)
	if err != nil {
		logger.Error("Failed to create feature flag",
			"error", err.Error(),
			"name", req.Name,
		)
		return errors.Internal("Failed to create feature flag")
	}

	return c.Status(201).JSON(flag)
}

// GetSystemBanners handles GET /admin/banners
func (h *AdminHandler) GetSystemBanners(c *fiber.Ctx) error {
	activeOnly := c.Query("active") == "true"

	banners, err := h.systemBannerRepo.GetAll(c.Context(), activeOnly)
	if err != nil {
		logger.Error("Failed to get system banners",
			"error", err.Error(),
		)
		return errors.Internal("Failed to get system banners")
	}

	return c.JSON(banners)
}

// CreateSystemBanner handles POST /admin/banners
func (h *AdminHandler) CreateSystemBanner(c *fiber.Ctx) error {
	var req models.CreateSystemBannerRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// TODO: Add validation using validator package
	if req.Message == "" {
		return errors.BadRequest("Banner message is required")
	}

	// Validate level
	if req.Level != models.BannerLevelInfo &&
		req.Level != models.BannerLevelWarning &&
		req.Level != models.BannerLevelError &&
		req.Level != models.BannerLevelSuccess {
		return errors.BadRequest("Invalid banner level. Must be one of: info, warning, error, success")
	}

	banner := &models.SystemBanner{
		Title:   req.Title,
		Message: req.Message,
		Level:   req.Level,
		Active:  req.Active,
	}

	err := h.systemBannerRepo.Create(c.Context(), banner)
	if err != nil {
		logger.Error("Failed to create system banner",
			"error", err.Error(),
			"message", req.Message,
		)
		return errors.Internal("Failed to create system banner")
	}

	return c.Status(201).JSON(banner)
}

// UpdateSystemBanner handles PUT /admin/banners/:id
func (h *AdminHandler) UpdateSystemBanner(c *fiber.Ctx) error {
	bannerIDStr := c.Params("id")
	bannerID, err := uuid.Parse(bannerIDStr)
	if err != nil {
		return errors.BadRequest("Invalid banner ID")
	}

	// Get existing banner
	banner, err := h.systemBannerRepo.GetByID(c.Context(), bannerID)
	if err != nil {
		logger.Error("Failed to get system banner",
			"error", err.Error(),
			"bannerID", bannerID,
		)
		if err.Error() == "system banner not found" {
			return errors.NotFound("System banner")
		}
		return errors.Internal("Failed to get system banner")
	}

	var req models.UpdateSystemBannerRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// Update fields if provided
	if req.Title != nil {
		banner.Title = req.Title
	}
	if req.Message != nil {
		banner.Message = *req.Message
	}
	if req.Level != nil {
		// Validate level
		if *req.Level != models.BannerLevelInfo &&
			*req.Level != models.BannerLevelWarning &&
			*req.Level != models.BannerLevelError &&
			*req.Level != models.BannerLevelSuccess {
			return errors.BadRequest("Invalid banner level. Must be one of: info, warning, error, success")
		}
		banner.Level = *req.Level
	}
	if req.Active != nil {
		banner.Active = *req.Active
	}

	err = h.systemBannerRepo.Update(c.Context(), banner)
	if err != nil {
		logger.Error("Failed to update system banner",
			"error", err.Error(),
			"bannerID", bannerID,
		)
		return errors.Internal("Failed to update system banner")
	}

	return c.JSON(banner)
}

// DeleteSystemBanner handles DELETE /admin/banners/:id
func (h *AdminHandler) DeleteSystemBanner(c *fiber.Ctx) error {
	bannerIDStr := c.Params("id")
	bannerID, err := uuid.Parse(bannerIDStr)
	if err != nil {
		return errors.BadRequest("Invalid banner ID")
	}

	err = h.systemBannerRepo.Delete(c.Context(), bannerID)
	if err != nil {
		logger.Error("Failed to delete system banner",
			"error", err.Error(),
			"bannerID", bannerID,
		)
		if err.Error() == "system banner not found" {
			return errors.NotFound("System banner")
		}
		return errors.Internal("Failed to delete system banner")
	}

	return c.SendStatus(204)
}