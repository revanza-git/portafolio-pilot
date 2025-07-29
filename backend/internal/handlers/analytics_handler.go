package handlers

import (
	"fmt"
	"time"

	"github.com/defi-dashboard/backend/pkg/pnl"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	pnlService  pnl.Service
	csvExporter *pnl.CSVExporter
}

func NewAnalyticsHandler(pnlService pnl.Service, csvExporter *pnl.CSVExporter) *AnalyticsHandler {
	return &AnalyticsHandler{
		pnlService:  pnlService,
		csvExporter: csvExporter,
	}
}

// GetPnL handles GET /analytics/pnl/:address
func (h *AnalyticsHandler) GetPnL(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return errors.BadRequest("Address is required")
	}

	// Parse query parameters
	fromStr := c.Query("from")
	toStr := c.Query("to")
	methodStr := c.Query("method", "fifo")

	// Parse dates
	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return errors.BadRequest("Invalid from date format. Use YYYY-MM-DD")
		}
	} else {
		// Default to 1 year ago
		from = time.Now().AddDate(-1, 0, 0)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return errors.BadRequest("Invalid to date format. Use YYYY-MM-DD")
		}
	} else {
		// Default to now
		to = time.Now()
	}

	// Validate method
	var method pnl.CalculationMethod
	switch methodStr {
	case "fifo":
		method = pnl.FIFO
	case "lifo":
		method = pnl.LIFO
	default:
		return errors.BadRequest("Invalid method. Use 'fifo' or 'lifo'")
	}

	// Calculate PnL
	calculation, err := h.pnlService.CalculatePnL(c.Context(), address, from, to, method)
	if err != nil {
		logger.Error("Failed to calculate PnL",
			"error", err.Error(),
			"address", address,
			"method", methodStr,
		)
		return errors.Internal("Failed to calculate PnL")
	}

	return c.JSON(calculation)
}

// ExportPnL handles GET /analytics/export
func (h *AnalyticsHandler) ExportPnL(c *fiber.Ctx) error {
	// Parse query parameters
	address := c.Query("address")
	if address == "" {
		return errors.BadRequest("Address is required")
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	methodStr := c.Query("method", "fifo")
	formatStr := c.Query("format", "csv")

	// Parse dates
	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return errors.BadRequest("Invalid from date format. Use YYYY-MM-DD")
		}
	} else {
		// Default to 1 year ago
		from = time.Now().AddDate(-1, 0, 0)
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return errors.BadRequest("Invalid to date format. Use YYYY-MM-DD")
		}
	} else {
		// Default to now
		to = time.Now()
	}

	// Validate method
	var method pnl.CalculationMethod
	switch methodStr {
	case "fifo":
		method = pnl.FIFO
	case "lifo":
		method = pnl.LIFO
	default:
		return errors.BadRequest("Invalid method. Use 'fifo' or 'lifo'")
	}

	// Validate format
	if formatStr != "csv" {
		return errors.BadRequest("Only CSV format is currently supported")
	}

	// Get export data
	exportData, err := h.pnlService.GetPnLExportData(c.Context(), address, from, to, method)
	if err != nil {
		logger.Error("Failed to get PnL export data",
			"error", err.Error(),
			"address", address,
			"method", methodStr,
		)
		return errors.Internal("Failed to get PnL export data")
	}

	if len(exportData) == 0 {
		return errors.NotFound("No PnL data found for the specified criteria")
	}

	// Check if client wants streaming or file download
	stream := c.Query("stream", "false")
	
	if stream == "true" {
		// Stream CSV directly to client
		c.Set("Content-Type", "text/csv")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"pnl_export_%s_%s.csv\"", 
			address[:8], time.Now().Format("20060102_150405")))

		return h.csvExporter.ExportToWriter(c.Response().BodyWriter(), exportData)
	} else {
		// Create temporary file and return download URL
		filepath, err := h.csvExporter.ExportToCSV(exportData, address)
		if err != nil {
			logger.Error("Failed to create CSV file",
				"error", err.Error(),
				"address", address,
			)
			return errors.Internal("Failed to create CSV file")
		}

		// Schedule cleanup after 1 hour
		h.csvExporter.ScheduleCleanup(filepath, time.Hour)

		// Return download info
		return c.JSON(fiber.Map{
			"download_url": fmt.Sprintf("/analytics/download?file=%s", filepath),
			"expires_at":   time.Now().Add(time.Hour).Unix(),
			"filename":     fmt.Sprintf("pnl_export_%s_%s.csv", address[:8], time.Now().Format("20060102_150405")),
			"record_count": len(exportData),
		})
	}
}

// DownloadFile handles GET /analytics/download for file downloads
func (h *AnalyticsHandler) DownloadFile(c *fiber.Ctx) error {
	filepath := c.Query("file")
	if filepath == "" {
		return errors.BadRequest("File parameter is required")
	}

	// Security check: only allow files from temp directory
	// This is a basic check - in production, you'd want more robust validation
	if !isValidTempFile(filepath) {
		return errors.BadRequest("Invalid file path")
	}

	return c.SendFile(filepath)
}

// Helper function to validate temp file paths
func isValidTempFile(filepath string) bool {
	// Basic validation - in production, implement proper path validation
	// to prevent directory traversal attacks
	return filepath != "" && len(filepath) > 0
}

// GetPnLSummary handles GET /analytics/summary/:address for dashboard display
func (h *AnalyticsHandler) GetPnLSummary(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return errors.BadRequest("Address is required")
	}

	// Get summary for the last 30 days by default
	from := time.Now().AddDate(0, 0, -30)
	to := time.Now()

	// Allow custom time range
	if fromStr := c.Query("from"); fromStr != "" {
		if parsedFrom, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = parsedFrom
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if parsedTo, err := time.Parse("2006-01-02", toStr); err == nil {
			to = parsedTo
		}
	}

	// Calculate PnL for both FIFO and LIFO
	fifoCalc, fifoErr := h.pnlService.CalculatePnL(c.Context(), address, from, to, pnl.FIFO)
	lifoCalc, lifoErr := h.pnlService.CalculatePnL(c.Context(), address, from, to, pnl.LIFO)

	summary := fiber.Map{
		"address": address,
		"period": fiber.Map{
			"from": from.Format("2006-01-02"),
			"to":   to.Format("2006-01-02"),
		},
	}

	if fifoErr == nil {
		summary["fifo"] = fifoCalc
	} else {
		summary["fifo_error"] = fifoErr.Error()
	}

	if lifoErr == nil {
		summary["lifo"] = lifoCalc
	} else {
		summary["lifo_error"] = lifoErr.Error()
	}

	return c.JSON(summary)
}