package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/defi-dashboard/backend/internal/models"
	"github.com/defi-dashboard/backend/pkg/logger"
	"github.com/google/uuid"
)

const (
	AlchemyMainnetURL = "https://eth-mainnet.g.alchemy.com/v2"
	AlchemyPolygonURL = "https://polygon-mainnet.g.alchemy.com/v2"
	AlchemyArbitrumURL = "https://arb-mainnet.g.alchemy.com/v2"
	AlchemyOptimismURL = "https://opt-mainnet.g.alchemy.com/v2"
)

type AlchemyClient struct {
	httpClient *http.Client
	apiKey     string
	baseURLs   map[int]string
}

func NewAlchemyClient(apiKey string) *AlchemyClient {
	return &AlchemyClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
		baseURLs: map[int]string{
			1:     fmt.Sprintf("%s/%s", AlchemyMainnetURL, apiKey),
			137:   fmt.Sprintf("%s/%s", AlchemyPolygonURL, apiKey),
			42161: fmt.Sprintf("%s/%s", AlchemyArbitrumURL, apiKey),
			10:    fmt.Sprintf("%s/%s", AlchemyOptimismURL, apiKey),
		},
	}
}

type TokenBalance struct {
	ContractAddress  string `json:"contractAddress"`
	TokenBalance     string `json:"tokenBalance"`
	Error            string `json:"error,omitempty"`
}

type TokenMetadata struct {
	Decimals int    `json:"decimals"`
	Logo     string `json:"logo"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
}

type AlchemyTokenBalanceResponse struct {
	Address      string         `json:"address"`
	TokenBalances []TokenBalance `json:"tokenBalances"`
}

type AlchemyTokenMetadataResponse struct {
	Data []TokenMetadata `json:"data"`
}

type AlchemyTransactionResponse struct {
	Result struct {
		Transfers []TransferData `json:"transfers"`
	} `json:"result"`
}

type TransferData struct {
	BlockNum    string        `json:"blockNum"`
	Hash        string        `json:"hash"`
	From        string        `json:"from"`
	To          string        `json:"to"`
	Value       float64       `json:"value"`
	Asset       string        `json:"asset"`
	Category    string        `json:"category"`
	RawContract RawContract   `json:"rawContract"`
	Metadata    TransferMeta  `json:"metadata"`
}

type RawContract struct {
	Value   string `json:"value"`
	Address string `json:"address"`
	Decimal string `json:"decimal"`
}

type TransferMeta struct {
	BlockTimestamp string `json:"blockTimestamp"`
}

// GetTokenBalances fetches ERC20 token balances for an address
func (c *AlchemyClient) GetTokenBalances(ctx context.Context, address string, chainID int) ([]*models.Balance, error) {
	baseURL, exists := c.baseURLs[chainID]
	if !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	// Get token balances
	reqBody := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_getTokenBalances",
		"params":  []interface{}{address},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var balanceResp struct {
		Result AlchemyTokenBalanceResponse `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&balanceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if balanceResp.Error != nil {
		return nil, fmt.Errorf("alchemy API error: %s", balanceResp.Error.Message)
	}

	// Get metadata for tokens with non-zero balances
	var tokenAddresses []string
	for _, balance := range balanceResp.Result.TokenBalances {
		if balance.TokenBalance != "0x0" && balance.Error == "" {
			tokenAddresses = append(tokenAddresses, balance.ContractAddress)
		}
	}

	if len(tokenAddresses) == 0 {
		return []*models.Balance{}, nil
	}

	// Get token metadata
	metadata, err := c.getTokenMetadata(ctx, tokenAddresses, chainID)
	if err != nil {
		logger.Error("Failed to get token metadata", "error", err)
		// Continue with empty metadata
		metadata = make(map[string]TokenMetadata)
	}

	// Convert to models.Balance
	var balances []*models.Balance
	for _, tokenBalance := range balanceResp.Result.TokenBalances {
		if tokenBalance.TokenBalance == "0x0" || tokenBalance.Error != "" {
			continue
		}

		meta, exists := metadata[tokenBalance.ContractAddress]
		if !exists {
			// Skip tokens without metadata
			continue
		}

		// Convert hex balance to decimal
		balanceInt := new(big.Int)
		balanceInt.SetString(tokenBalance.TokenBalance[2:], 16) // Remove 0x prefix

		// Create token
		token := &models.Token{
			ID:       uuid.New(),
			Address:  tokenBalance.ContractAddress,
			ChainID:  chainID,
			Symbol:   meta.Symbol,
			Name:     meta.Name,
			Decimals: meta.Decimals,
			LogoURI:  &meta.Logo,
		}

		balance := &models.Balance{
			ID:       uuid.New(),
			WalletID: uuid.New(), // This should be set by the service
			TokenID:  token.ID,
			Token:    token,
			Balance:  balanceInt.String(),
		}

		balances = append(balances, balance)
	}

	return balances, nil
}

// getTokenMetadata fetches metadata for multiple tokens
func (c *AlchemyClient) getTokenMetadata(ctx context.Context, addresses []string, chainID int) (map[string]TokenMetadata, error) {
	baseURL, exists := c.baseURLs[chainID]
	if !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	reqBody := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_getTokenMetadata",
		"params":  []interface{}{addresses},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var metadataResp struct {
		Result AlchemyTokenMetadataResponse `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&metadataResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if metadataResp.Error != nil {
		return nil, fmt.Errorf("alchemy API error: %s", metadataResp.Error.Message)
	}

	// Create address -> metadata mapping
	result := make(map[string]TokenMetadata)
	for i, meta := range metadataResp.Result.Data {
		if i < len(addresses) {
			result[addresses[i]] = meta
		}
	}

	return result, nil
}

// GetETHBalance fetches native ETH balance for an address
func (c *AlchemyClient) GetETHBalance(ctx context.Context, address string, chainID int) (*big.Int, error) {
	baseURL, exists := c.baseURLs[chainID]
	if !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	reqBody := map[string]interface{}{
		"id":      1,  
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  []interface{}{address, "latest"},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var balanceResp struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&balanceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if balanceResp.Error != nil {
		return nil, fmt.Errorf("alchemy API error: %s", balanceResp.Error.Message)
	}

	// Convert hex to big.Int
	balance := new(big.Int)
	balance.SetString(balanceResp.Result[2:], 16) // Remove 0x prefix

	return balance, nil
}

// GetTransactions fetches recent transactions for an address
func (c *AlchemyClient) GetTransactions(ctx context.Context, address string, chainID int) ([]*models.Transaction, error) {
	baseURL, exists := c.baseURLs[chainID]
	if !exists {
		return nil, fmt.Errorf("unsupported chain ID: %d", chainID)
	}

	reqBody := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_getAssetTransfers",
		"params": []map[string]interface{}{
			{
				"fromBlock":         "0x0",
				"toBlock":           "latest",
				"fromAddress":       address,
				"category":          []string{"external", "internal", "erc20", "erc721", "erc1155"},
				"withMetadata":      true,
				"excludeZeroValue":  true,
				"maxCount":          "0x64", // 100 transactions
			},
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var txResp AlchemyTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&txResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to models.Transaction
	var transactions []*models.Transaction
	for _, transfer := range txResp.Result.Transfers {
		blockNum, _ := strconv.ParseInt(transfer.BlockNum[2:], 16, 64)
		
		// Parse timestamp
		timestamp, err := time.Parse(time.RFC3339, transfer.Metadata.BlockTimestamp)
		if err != nil {
			timestamp = time.Now() // Fallback
		}

		tx := &models.Transaction{
			ID:          uuid.New(),
			Hash:        transfer.Hash,
			ChainID:     chainID,
			FromAddress: transfer.From,
			ToAddress:   &transfer.To,
			Value:       &transfer.RawContract.Value,
			BlockNumber: &blockNum,
			Timestamp:   timestamp,
			Status:      "success", // Alchemy only returns successful transfers
			Type:        transfer.Category,
			Metadata: map[string]interface{}{
				"asset":    transfer.Asset,
				"category": transfer.Category,
			},
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}