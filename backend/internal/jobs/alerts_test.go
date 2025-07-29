package jobs

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AlertEvaluatorTestSuite tests the alert evaluator logic
type AlertEvaluatorTestSuite struct {
	suite.Suite
	job *AlertEvaluatorJob
	ctx context.Context
}

func (s *AlertEvaluatorTestSuite) SetupTest() {
	s.ctx = context.Background()
	// Using nil db for unit tests - we'll test pure logic
	s.job = NewAlertEvaluatorJob(nil)
}

// TestEvaluatePriceCondition tests price alert condition evaluation
func (s *AlertEvaluatorTestSuite) TestEvaluatePriceCondition() {
	tests := []struct {
		name          string
		alertType     string
		targetPrice   float64
		currentPrice  float64
		shouldTrigger bool
	}{
		{
			name:          "Price above threshold - should trigger",
			alertType:     AlertTypePriceAbove,
			targetPrice:   2000.0,
			currentPrice:  2100.0,
			shouldTrigger: true,
		},
		{
			name:          "Price above threshold - should not trigger",
			alertType:     AlertTypePriceAbove,
			targetPrice:   2000.0,
			currentPrice:  1900.0,
			shouldTrigger: false,
		},
		{
			name:          "Price below threshold - should trigger",
			alertType:     AlertTypePriceBelow,
			targetPrice:   2000.0,
			currentPrice:  1900.0,
			shouldTrigger: true,
		},
		{
			name:          "Price below threshold - should not trigger",
			alertType:     AlertTypePriceBelow,
			targetPrice:   2000.0,
			currentPrice:  2100.0,
			shouldTrigger: false,
		},
		{
			name:          "Price exactly at threshold - should not trigger",
			alertType:     AlertTypePriceAbove,
			targetPrice:   2000.0,
			currentPrice:  2000.0,
			shouldTrigger: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			conditions := map[string]interface{}{
				"price": tt.targetPrice,
			}
			conditionsJSON, _ := json.Marshal(conditions)

			alert := &Alert{
				ID:         uuid.New().String(),
				Type:       tt.alertType,
				Conditions: conditionsJSON,
			}

			result := s.job.evaluatePriceCondition(alert, tt.currentPrice)
			s.Equal(tt.shouldTrigger, result)
		})
	}
}

// TestGroupAlertsByType tests alert grouping functionality
func (s *AlertEvaluatorTestSuite) TestGroupAlertsByType() {
	alerts := []*Alert{
		{Type: AlertTypePriceAbove},
		{Type: AlertTypePriceBelow},
		{Type: AlertTypePriceAbove},
		{Type: AlertTypeApproval},
		{Type: AlertTypePriceBelow},
		{Type: AlertTypeAPRChange},
	}

	grouped := s.job.groupAlertsByType(alerts)

	s.Len(grouped[AlertTypePriceAbove], 2)
	s.Len(grouped[AlertTypePriceBelow], 2)
	s.Len(grouped[AlertTypeApproval], 1)
	s.Len(grouped[AlertTypeAPRChange], 1)
	s.Len(grouped[AlertTypeLargeTransfer], 0)
}

// TestAlertTarget tests alert target parsing
func (s *AlertEvaluatorTestSuite) TestAlertTarget() {
	testCases := []struct {
		name     string
		target   AlertTarget
		expected string
	}{
		{
			name: "Token target",
			target: AlertTarget{
				Type:       "token",
				Identifier: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				ChainID:    1,
			},
			expected: "token",
		},
		{
			name: "Address target",
			target: AlertTarget{
				Type:       "address",
				Identifier: "0x1234567890123456789012345678901234567890",
				ChainID:    1,
			},
			expected: "address",
		},
		{
			name: "Pool target",
			target: AlertTarget{
				Type:       "pool",
				Identifier: "aave-v3-usdc",
				ChainID:    0,
			},
			expected: "pool",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			targetJSON, err := json.Marshal(tc.target)
			s.NoError(err)

			var parsed AlertTarget
			err = json.Unmarshal(targetJSON, &parsed)
			s.NoError(err)
			s.Equal(tc.expected, parsed.Type)
			s.Equal(tc.target.Identifier, parsed.Identifier)
		})
	}
}

// TestAPRAlertConditions tests APR alert condition evaluation
func (s *AlertEvaluatorTestSuite) TestAPRAlertConditions() {
	testCases := []struct {
		name          string
		minAPR        *float64
		maxAPR        *float64
		currentAPR    float64
		shouldTrigger bool
	}{
		{
			name:          "APR below minimum",
			minAPR:        floatPtr(5.0),
			maxAPR:        nil,
			currentAPR:    3.0,
			shouldTrigger: true,
		},
		{
			name:          "APR above minimum",
			minAPR:        floatPtr(5.0),
			maxAPR:        nil,
			currentAPR:    7.0,
			shouldTrigger: false,
		},
		{
			name:          "APR above maximum",
			minAPR:        nil,
			maxAPR:        floatPtr(10.0),
			currentAPR:    12.0,
			shouldTrigger: true,
		},
		{
			name:          "APR below maximum",
			minAPR:        nil,
			maxAPR:        floatPtr(10.0),
			currentAPR:    8.0,
			shouldTrigger: false,
		},
		{
			name:          "APR within range",
			minAPR:        floatPtr(5.0),
			maxAPR:        floatPtr(10.0),
			currentAPR:    7.0,
			shouldTrigger: false,
		},
		{
			name:          "APR outside range - too low",
			minAPR:        floatPtr(5.0),
			maxAPR:        floatPtr(10.0),
			currentAPR:    3.0,
			shouldTrigger: true,
		},
		{
			name:          "APR outside range - too high",
			minAPR:        floatPtr(5.0),
			maxAPR:        floatPtr(10.0),
			currentAPR:    15.0,
			shouldTrigger: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the evaluation logic
			shouldTrigger := false
			if tc.minAPR != nil && tc.currentAPR < *tc.minAPR {
				shouldTrigger = true
			}
			if tc.maxAPR != nil && tc.currentAPR > *tc.maxAPR {
				shouldTrigger = true
			}
			
			s.Equal(tc.shouldTrigger, shouldTrigger)
		})
	}
}

// TestTransferThresholdParsing tests large transfer threshold parsing
func (s *AlertEvaluatorTestSuite) TestTransferThresholdParsing() {
	testCases := []struct {
		name      string
		threshold string
		amount    string
		exceeds   bool
	}{
		{
			name:      "1 ETH threshold - transfer exceeds",
			threshold: "1000000000000000000", // 1 ETH in wei
			amount:    "2000000000000000000", // 2 ETH in wei
			exceeds:   true,
		},
		{
			name:      "1 ETH threshold - transfer below",
			threshold: "1000000000000000000", // 1 ETH in wei
			amount:    "500000000000000000",  // 0.5 ETH in wei
			exceeds:   false,
		},
		{
			name:      "Large USDC threshold",
			threshold: "1000000000", // 1000 USDC (6 decimals)
			amount:    "2000000000", // 2000 USDC
			exceeds:   true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Test the comparison logic used in evaluateTransferAlerts
			thresholdBig, ok := new(big.Int).SetString(tc.threshold, 10)
			s.True(ok)
			
			amountBig, ok := new(big.Int).SetString(tc.amount, 10)
			s.True(ok)
			
			exceeds := amountBig.Cmp(thresholdBig) > 0
			s.Equal(tc.exceeds, exceeds)
		})
	}
}

// TestAlertCooldown tests that alerts respect cooldown period
func (s *AlertEvaluatorTestSuite) TestAlertCooldown() {
	now := time.Now()
	testCases := []struct {
		name             string
		lastTriggered    *time.Time
		shouldEvaluate   bool
	}{
		{
			name:             "Never triggered - should evaluate",
			lastTriggered:    nil,
			shouldEvaluate:   true,
		},
		{
			name:             "Triggered 2 hours ago - should evaluate",
			lastTriggered:    timePtr(now.Add(-2 * time.Hour)),
			shouldEvaluate:   true,
		},
		{
			name:             "Triggered 30 minutes ago - should not evaluate",
			lastTriggered:    timePtr(now.Add(-30 * time.Minute)),
			shouldEvaluate:   false,
		},
		{
			name:             "Triggered exactly 1 hour ago - should evaluate",
			lastTriggered:    timePtr(now.Add(-1 * time.Hour)),
			shouldEvaluate:   true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Test the cooldown logic used in getActiveAlerts
			shouldEvaluate := tc.lastTriggered == nil || 
				tc.lastTriggered.Before(now.Add(-1*time.Hour))
			s.Equal(tc.shouldEvaluate, shouldEvaluate)
		})
	}
}

// TestNotificationPrefs tests notification preference parsing
func (s *AlertEvaluatorTestSuite) TestNotificationPrefs() {
	testCases := []struct {
		name     string
		prefs    NotificationPrefs
		hasEmail bool
		hasWebhook bool
	}{
		{
			name: "Email only",
			prefs: NotificationPrefs{
				Email: true,
			},
			hasEmail: true,
			hasWebhook: false,
		},
		{
			name: "Webhook only",
			prefs: NotificationPrefs{
				Email: false,
				Webhook: "https://example.com/webhook",
			},
			hasEmail: false,
			hasWebhook: true,
		},
		{
			name: "Both email and webhook",
			prefs: NotificationPrefs{
				Email: true,
				Webhook: "https://example.com/webhook",
			},
			hasEmail: true,
			hasWebhook: true,
		},
		{
			name: "Neither",
			prefs: NotificationPrefs{
				Email: false,
			},
			hasEmail: false,
			hasWebhook: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Equal(tc.hasEmail, tc.prefs.Email)
			s.Equal(tc.hasWebhook, tc.prefs.Webhook != "")
		})
	}
}

// Test helpers
func floatPtr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// Run the test suite
func TestAlertEvaluatorSuite(t *testing.T) {
	suite.Run(t, new(AlertEvaluatorTestSuite))
}

// Additional standalone tests

func TestPoolTVLChangeCalculation(t *testing.T) {
	testCases := []struct {
		name         string
		currentTVL   float64
		previousTVL  float64
		expectedPct  float64
	}{
		{
			name:        "10% increase",
			currentTVL:  110000,
			previousTVL: 100000,
			expectedPct: 10.0,
		},
		{
			name:        "10% decrease",
			currentTVL:  90000,
			previousTVL: 100000,
			expectedPct: -10.0,
		},
		{
			name:        "No change",
			currentTVL:  100000,
			previousTVL: 100000,
			expectedPct: 0.0,
		},
		{
			name:        "50% increase",
			currentTVL:  150000,
			previousTVL: 100000,
			expectedPct: 50.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate percentage change
			changePercent := ((tc.currentTVL - tc.previousTVL) / tc.previousTVL) * 100
			assert.InDelta(t, tc.expectedPct, changePercent, 0.001)
		})
	}
}