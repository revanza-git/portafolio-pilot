package handlers

import (
	"strconv"

	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// GetTransactions handles GET /transactions/:address
func (h *TransactionHandler) GetTransactions(c *fiber.Ctx) error {
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

	var txType *string
	if typeParam := c.Query("type"); typeParam != "" {
		// Validate transaction type
		validTypes := map[string]bool{
			"send": true, "receive": true, "swap": true,
			"approve": true, "bridge": true, "stake": true, "unstake": true,
		}
		if !validTypes[typeParam] {
			return errors.BadRequest("Invalid transaction type")
		}
		txType = &typeParam
	}

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Get transactions
	transactions, err := h.transactionService.GetTransactions(c.Context(), address, chainID, txType, page, limit)
	if err != nil {
		return err
	}

	return c.JSON(transactions)
}

// GetApprovals handles GET /transactions/:address/approvals
func (h *TransactionHandler) GetApprovals(c *fiber.Ctx) error {
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

	activeOnly := c.Query("active", "true") == "true"

	// Get approvals
	approvals, err := h.transactionService.GetApprovals(c.Context(), address, chainID, activeOnly)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"approvals": approvals,
	})
}

// RevokeApproval handles DELETE /transactions/:address/approvals/:token
func (h *TransactionHandler) RevokeApproval(c *fiber.Ctx) error {
	address := c.Params("address")
	token := c.Params("token")
	spender := c.Query("spender")

	if address == "" || token == "" || spender == "" {
		return errors.BadRequest("Address, token, and spender are required")
	}

	// TODO: Validate that the authenticated user owns this address
	authAddress := c.Locals("address").(string)
	if authAddress != address {
		return errors.Forbidden("You can only revoke approvals for your own address")
	}

	// Revoke approval
	txHash, err := h.transactionService.RevokeApproval(c.Context(), address, token, spender)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"txHash": txHash,
	})
}