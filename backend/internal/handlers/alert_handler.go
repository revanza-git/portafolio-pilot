package handlers

import (
	"strconv"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AlertHandler struct {
	alertService services.AlertService
}

func NewAlertHandler(alertService services.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

// CreateAlert handles POST /alerts
func (h *AlertHandler) CreateAlert(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	var req models.CreateAlertRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	// TODO: Add validation using validator package
	if req.Type == "" {
		return errors.BadRequest("Alert type is required")
	}

	alert, err := h.alertService.CreateAlert(c.Context(), userID, &req)
	if err != nil {
		logger.Error("Failed to create alert",
			"error", err.Error(),
			"userID", userID,
			"type", req.Type,
		)
		return errors.Internal("Failed to create alert")
	}

	return c.Status(201).JSON(alert)
}

// GetAlerts handles GET /alerts
func (h *AlertHandler) GetAlerts(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	// Parse query parameters
	var status *string
	if statusParam := c.Query("status"); statusParam != "" {
		status = &statusParam
	}

	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	alerts, err := h.alertService.GetUserAlerts(c.Context(), userID, status, limit, offset)
	if err != nil {
		logger.Error("Failed to get alerts",
			"error", err.Error(),
			"userID", userID,
		)
		return errors.Internal("Failed to get alerts")
	}

	// Calculate total for pagination (simplified - in production you'd want a separate count query)
	total := len(alerts)
	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"data": alerts,
		"meta": fiber.Map{
			"page":       (offset / limit) + 1,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetAlert handles GET /alerts/:alertId
func (h *AlertHandler) GetAlert(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	alertIDStr := c.Params("alertId")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		return errors.BadRequest("Invalid alert ID")
	}

	alert, err := h.alertService.GetAlert(c.Context(), alertID, userID)
	if err != nil {
		logger.Error("Failed to get alert",
			"error", err.Error(),
			"alertID", alertID,
			"userID", userID,
		)
		return errors.NotFound("Alert")
	}

	return c.JSON(alert)
}

// UpdateAlert handles PATCH /alerts/:alertId
func (h *AlertHandler) UpdateAlert(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	alertIDStr := c.Params("alertId")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		return errors.BadRequest("Invalid alert ID")
	}

	var req models.UpdateAlertRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	alert, err := h.alertService.UpdateAlert(c.Context(), alertID, userID, &req)
	if err != nil {
		logger.Error("Failed to update alert",
			"error", err.Error(),
			"alertID", alertID,
			"userID", userID,
		)
		if err.Error() == "alert not found" {
			return errors.NotFound("Alert")
		}
		return errors.Internal("Failed to update alert")
	}

	return c.JSON(alert)
}

// DeleteAlert handles DELETE /alerts/:alertId
func (h *AlertHandler) DeleteAlert(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	alertIDStr := c.Params("alertId")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		return errors.BadRequest("Invalid alert ID")
	}

	err = h.alertService.DeleteAlert(c.Context(), alertID, userID)
	if err != nil {
		logger.Error("Failed to delete alert",
			"error", err.Error(),
			"alertID", alertID,
			"userID", userID,
		)
		if err.Error() == "alert not found" {
			return errors.NotFound("Alert")
		}
		return errors.Internal("Failed to delete alert")
	}

	return c.SendStatus(204)
}

// GetAlertHistory handles GET /alerts/history
func (h *AlertHandler) GetAlertHistory(c *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	// Parse optional alert ID
	var alertID *uuid.UUID
	if alertIDStr := c.Query("alertId"); alertIDStr != "" {
		id, err := uuid.Parse(alertIDStr)
		if err != nil {
			return errors.BadRequest("Invalid alert ID")
		}
		alertID = &id
	}

	// Parse pagination parameters
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	history, err := h.alertService.GetAlertHistory(c.Context(), alertID, userID, limit, offset)
	if err != nil {
		logger.Error("Failed to get alert history",
			"error", err.Error(),
			"alertID", alertID,
			"userID", userID,
		)
		return errors.Internal("Failed to get alert history")
	}

	// Calculate total for pagination
	total := len(history)
	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"data": history,
		"meta": fiber.Map{
			"page":       (offset / limit) + 1,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// PauseAlert handles PATCH /alerts/:alertId/pause
func (h *AlertHandler) PauseAlert(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	alertIDStr := c.Params("alertId")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		return errors.BadRequest("Invalid alert ID")
	}

	status := models.AlertStatusDisabled
	req := models.UpdateAlertRequest{
		Status: &status,
	}

	alert, err := h.alertService.UpdateAlert(c.Context(), alertID, userID, &req)
	if err != nil {
		logger.Error("Failed to pause alert",
			"error", err.Error(),
			"alertID", alertID,
			"userID", userID,
		)
		if err.Error() == "alert not found" {
			return errors.NotFound("Alert")
		}
		return errors.Internal("Failed to pause alert")
	}

	return c.JSON(alert)
}

// ActivateAlert handles PATCH /alerts/:alertId/activate
func (h *AlertHandler) ActivateAlert(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return errors.Unauthorized("User not authenticated")
	}

	alertIDStr := c.Params("alertId")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		return errors.BadRequest("Invalid alert ID")
	}

	status := models.AlertStatusActive
	req := models.UpdateAlertRequest{
		Status: &status,
	}

	alert, err := h.alertService.UpdateAlert(c.Context(), alertID, userID, &req)
	if err != nil {
		logger.Error("Failed to activate alert",
			"error", err.Error(),
			"alertID", alertID,
			"userID", userID,
		)
		if err.Error() == "alert not found" {
			return errors.NotFound("Alert")
		}
		return errors.Internal("Failed to activate alert")
	}

	return c.JSON(alert)
}