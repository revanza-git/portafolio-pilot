package jobs

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/internal/repos"
	"github.com/defi-dashboard/backend/internal/services"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertEvaluatorJob struct {
	db           *pgxpool.Pool
	alertService services.AlertService
	alertRepo    repos.AlertRepository
}

func NewAlertEvaluatorJob(db *pgxpool.Pool, alertService services.AlertService, alertRepo repos.AlertRepository) *AlertEvaluatorJob {
	return &AlertEvaluatorJob{
		db:           db,
		alertService: alertService,
		alertRepo:    alertRepo,
	}
}

// Use alert types from models
const (
	AlertTypePriceAbove      = models.AlertTypePriceAbove
	AlertTypePriceBelow      = models.AlertTypePriceBelow
	AlertTypeLargeTransfer   = models.AlertTypeLargeTransfer
	AlertTypeApproval        = models.AlertTypeApproval
	AlertTypeLiquidityChange = models.AlertTypeLiquidityChange
	AlertTypeAPRChange       = models.AlertTypeAPRChange
)

// Run executes the alert evaluation job
func (j *AlertEvaluatorJob) Run(ctx context.Context) error {
	logger.Info("Starting alert evaluation job")

	// Get active alerts
	alerts, err := j.getActiveAlerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active alerts: %w", err)
	}

	if len(alerts) == 0 {
		logger.Info("No active alerts to evaluate")
		return nil
	}

	logger.Info("Evaluating alerts", "count", len(alerts))

	// Group alerts by type for batch processing
	alertsByType := j.groupAlertsByType(alerts)

	// Evaluate each type of alert
	triggered := 0
	for alertType, typeAlerts := range alertsByType {
		count, err := j.evaluateAlertType(ctx, alertType, typeAlerts)
		if err != nil {
			logger.Error("Failed to evaluate alert type",
				"type", alertType,
				"error", err)
			continue
		}
		triggered += count
	}

	logger.Info("Alert evaluation completed",
		"total", len(alerts),
		"triggered", triggered)

	return nil
}

// getActiveAlerts retrieves all active alerts from the database
func (j *AlertEvaluatorJob) getActiveAlerts(ctx context.Context) ([]models.Alert, error) {
	return j.alertRepo.GetActiveAlerts(ctx)
}

// groupAlertsByType groups alerts by their type for batch processing
func (j *AlertEvaluatorJob) groupAlertsByType(alerts []models.Alert) map[string][]models.Alert {
	grouped := make(map[string][]models.Alert)
	for _, alert := range alerts {
		grouped[alert.Type] = append(grouped[alert.Type], alert)
	}
	return grouped
}

// evaluateAlertType evaluates all alerts of a specific type
func (j *AlertEvaluatorJob) evaluateAlertType(ctx context.Context, alertType string, alerts []models.Alert) (int, error) {
	switch alertType {
	case AlertTypePriceAbove, AlertTypePriceBelow:
		return j.evaluatePriceAlerts(ctx, alerts)
	case AlertTypeLargeTransfer:
		return j.evaluateTransferAlerts(ctx, alerts)
	case AlertTypeApproval:
		return j.evaluateApprovalAlerts(ctx, alerts)
	case AlertTypeLiquidityChange:
		return j.evaluateLiquidityAlerts(ctx, alerts)
	case AlertTypeAPRChange:
		return j.evaluateAPRAlerts(ctx, alerts)
	default:
		logger.Warn("Unknown alert type", "type", alertType)
		return 0, nil
	}
}

// evaluatePriceAlerts checks price-based alerts
func (j *AlertEvaluatorJob) evaluatePriceAlerts(ctx context.Context, alerts []models.Alert) (int, error) {
	// Get unique tokens to check
	tokenMap := make(map[string][]models.Alert)
	for _, alert := range alerts {
		if alert.Target.Type == "token" {
			key := fmt.Sprintf("%s-%d", alert.Target.Identifier, alert.Target.ChainID)
			tokenMap[key] = append(tokenMap[key], alert)
		}
	}

	// Fetch current prices
	prices, err := j.getTokenPrices(ctx, tokenMap)
	if err != nil {
		return 0, fmt.Errorf("failed to get token prices: %w", err)
	}

	// Evaluate each alert
	triggered := 0
	for tokenKey, tokenAlerts := range tokenMap {
		price, exists := prices[tokenKey]
		if !exists {
			continue
		}

		for _, alert := range tokenAlerts {
			if j.evaluatePriceCondition(&alert, price) {
				triggeredValue := map[string]interface{}{
					"currentPrice": price,
					"tokenKey":     tokenKey,
				}
				
				if err := j.alertService.TriggerAlert(ctx, alert.ID, triggeredValue); err != nil {
					logger.Error("Failed to trigger alert",
						"alertId", alert.ID,
						"error", err)
					continue
				}
				triggered++
			}
		}
	}

	return triggered, nil
}

// evaluatePriceCondition checks if a price alert condition is met
func (j *AlertEvaluatorJob) evaluatePriceCondition(alert *models.Alert, currentPrice float64) bool {
	if alert.Conditions.Price == nil {
		logger.Error("Price condition missing for price alert", "alertId", alert.ID)
		return false
	}

	targetPrice := *alert.Conditions.Price

	switch alert.Type {
	case AlertTypePriceAbove:
		return currentPrice > targetPrice
	case AlertTypePriceBelow:
		return currentPrice < targetPrice
	default:
		return false
	}
}

// evaluateTransferAlerts checks for large transfers
func (j *AlertEvaluatorJob) evaluateTransferAlerts(ctx context.Context, alerts []models.Alert) (int, error) {
	// Group by address
	addressMap := make(map[string][]models.Alert)
	for _, alert := range alerts {
		if alert.Target.Type == "address" {
			addressMap[alert.Target.Identifier] = append(addressMap[alert.Target.Identifier], alert)
		}
	}

	triggered := 0
	for address, addrAlerts := range addressMap {
		// Get recent large transfers
		transfers, err := j.getLargeTransfers(ctx, address)
		if err != nil {
			logger.Error("Failed to get transfers",
				"address", address,
				"error", err)
			continue
		}

		for _, alert := range addrAlerts {
			if alert.Conditions.Threshold == nil {
				continue
			}

			threshold, ok := new(big.Int).SetString(*alert.Conditions.Threshold, 10)
			if !ok {
				continue
			}

			// Check if any transfer exceeds threshold
			for _, transfer := range transfers {
				amount, ok := new(big.Int).SetString(transfer.Amount, 10)
				if !ok {
					continue
				}

				if amount.Cmp(threshold) > 0 {
					triggeredValue := map[string]interface{}{
						"transferAmount": transfer.Amount,
						"threshold":      *alert.Conditions.Threshold,
						"address":        address,
					}
					
					if err := j.alertService.TriggerAlert(ctx, alert.ID, triggeredValue); err != nil {
						logger.Error("Failed to trigger alert",
							"alertId", alert.ID,
							"error", err)
					} else {
						triggered++
					}
					break // Only trigger once per alert
				}
			}
		}
	}

	return triggered, nil
}

// evaluateApprovalAlerts checks for new token approvals
func (j *AlertEvaluatorJob) evaluateApprovalAlerts(ctx context.Context, alerts []models.Alert) (int, error) {
	triggered := 0
	
	for _, alert := range alerts {
		if alert.Target.Type != "address" {
			continue
		}

		// Check for new approvals since last trigger
		newApprovals, err := j.getNewApprovals(ctx, alert.Target.Identifier, alert.LastTriggeredAt)
		if err != nil {
			logger.Error("Failed to get approvals",
				"address", alert.Target.Identifier,
				"error", err)
			continue
		}

		if newApprovals > 0 {
			triggeredValue := map[string]interface{}{
				"newApprovals": newApprovals,
				"address":      alert.Target.Identifier,
			}
			
			if err := j.alertService.TriggerAlert(ctx, alert.ID, triggeredValue); err != nil {
				logger.Error("Failed to trigger alert",
					"alertId", alert.ID,
					"error", err)
			} else {
				triggered++
			}
		}
	}

	return triggered, nil
}

// evaluateLiquidityAlerts checks for liquidity changes in pools
func (j *AlertEvaluatorJob) evaluateLiquidityAlerts(ctx context.Context, alerts []models.Alert) (int, error) {
	triggered := 0
	
	for _, alert := range alerts {
		if alert.Target.Type != "pool" {
			continue
		}

		if alert.Conditions.ChangePercent == nil {
			continue
		}

		changeThreshold := *alert.Conditions.ChangePercent

		// Get pool TVL change
		tvlChange, err := j.getPoolTVLChange(ctx, alert.Target.Identifier)
		if err != nil {
			logger.Error("Failed to get TVL change",
				"pool", alert.Target.Identifier,
				"error", err)
			continue
		}

		if tvlChange > changeThreshold || tvlChange < -changeThreshold {
			triggeredValue := map[string]interface{}{
				"tvlChangePercent": tvlChange,
				"threshold":        changeThreshold,
				"poolId":           alert.Target.Identifier,
			}
			
			if err := j.alertService.TriggerAlert(ctx, alert.ID, triggeredValue); err != nil {
				logger.Error("Failed to trigger alert",
					"alertId", alert.ID,
					"error", err)
			} else {
				triggered++
			}
		}
	}

	return triggered, nil
}

// evaluateAPRAlerts checks for APR changes in yield pools
func (j *AlertEvaluatorJob) evaluateAPRAlerts(ctx context.Context, alerts []models.Alert) (int, error) {
	triggered := 0
	
	for _, alert := range alerts {
		if alert.Target.Type != "pool" {
			continue
		}

		// Get current pool APR
		currentAPR, err := j.getPoolAPR(ctx, alert.Target.Identifier)
		if err != nil {
			logger.Error("Failed to get pool APR",
				"pool", alert.Target.Identifier,
				"error", err)
			continue
		}

		shouldTrigger := false
		var triggerReason string
		
		if alert.Conditions.MinAPR != nil && currentAPR < *alert.Conditions.MinAPR {
			shouldTrigger = true
			triggerReason = "below_min_apr"
		}
		if alert.Conditions.MaxAPR != nil && currentAPR > *alert.Conditions.MaxAPR {
			shouldTrigger = true
			triggerReason = "above_max_apr"
		}

		if shouldTrigger {
			triggeredValue := map[string]interface{}{
				"currentAPR": currentAPR,
				"minAPR":     alert.Conditions.MinAPR,
				"maxAPR":     alert.Conditions.MaxAPR,
				"reason":     triggerReason,
				"poolId":     alert.Target.Identifier,
			}
			
			if err := j.alertService.TriggerAlert(ctx, alert.ID, triggeredValue); err != nil {
				logger.Error("Failed to trigger alert",
					"alertId", alert.ID,
					"error", err)
			} else {
				triggered++
			}
		}
	}

	return triggered, nil
}

// Helper methods to fetch data

func (j *AlertEvaluatorJob) getTokenPrices(ctx context.Context, tokenMap map[string][]models.Alert) (map[string]float64, error) {
	prices := make(map[string]float64)
	
	for tokenKey := range tokenMap {
		// Parse token key (format: "address-chainId")
		var address string
		var chainID int
		fmt.Sscanf(tokenKey, "%s-%d", &address, &chainID)

		var price float64
		err := j.db.QueryRow(ctx, `
			SELECT price_usd 
			FROM tokens 
			WHERE address = $1 AND chain_id = $2`,
			address, chainID).Scan(&price)
		
		if err == nil {
			prices[tokenKey] = price
		}
	}

	return prices, nil
}

type Transfer struct {
	Amount string
}

func (j *AlertEvaluatorJob) getLargeTransfers(ctx context.Context, address string) ([]Transfer, error) {
	rows, err := j.db.Query(ctx, `
		SELECT value 
		FROM transactions 
		WHERE (from_address = $1 OR to_address = $1)
			AND timestamp > NOW() - INTERVAL '1 hour'
			AND status = 'confirmed'
		ORDER BY timestamp DESC`,
		address)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var t Transfer
		if err := rows.Scan(&t.Amount); err == nil && t.Amount != "" {
			transfers = append(transfers, t)
		}
	}

	return transfers, rows.Err()
}

func (j *AlertEvaluatorJob) getNewApprovals(ctx context.Context, address string, since *time.Time) (int, error) {
	sinceTime := time.Now().Add(-1 * time.Hour)
	if since != nil {
		sinceTime = *since
	}

	var count int
	err := j.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM token_allowances ta
		INNER JOIN wallets w ON w.id = ta.wallet_id
		WHERE w.address = $1 
			AND ta.created_at > $2`,
		address, sinceTime).Scan(&count)
	
	return count, err
}

func (j *AlertEvaluatorJob) getPoolTVLChange(ctx context.Context, poolID string) (float64, error) {
	var currentTVL, previousTVL float64
	
	err := j.db.QueryRow(ctx, `
		SELECT tvl_usd,
			   LAG(tvl_usd) OVER (ORDER BY updated_at DESC) as prev_tvl
		FROM yield_pools
		WHERE pool_id = $1
		ORDER BY updated_at DESC
		LIMIT 1`,
		poolID).Scan(&currentTVL, &previousTVL)
	
	if err != nil || previousTVL == 0 {
		return 0, err
	}

	changePercent := ((currentTVL - previousTVL) / previousTVL) * 100
	return changePercent, nil
}

func (j *AlertEvaluatorJob) getPoolAPR(ctx context.Context, poolID string) (float64, error) {
	var apr float64
	err := j.db.QueryRow(ctx, `
		SELECT apy 
		FROM yield_pools 
		WHERE pool_id = $1`,
		poolID).Scan(&apr)
	
	return apr, err
}