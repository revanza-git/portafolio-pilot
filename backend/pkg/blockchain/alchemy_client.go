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
	PolygonAmoyURL = "https://rpc-amoy.polygon.technology" // Public RPC for testnet
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
			80002: PolygonAmoyURL, // Polygon Amoy testnet uses public RPC, no API key needed
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

	// Special handling for Polygon Amoy (public RPC, no Alchemy methods)
	if chainID == 80002 {
		return c.getTokenBalancesPublicRPC(ctx, address, chainID, baseURL)
	}

	// Get token balances using Alchemy-specific method
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

	// Special handling for Polygon Amoy (public RPC, no Alchemy methods)
	if chainID == 80002 {
		return c.getTransactionsPublicRPC(ctx, address, chainID, baseURL)
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

// Common ERC20 token addresses on Polygon Amoy testnet
var polygonAmoyTokens = map[string]struct {
	Symbol   string
	Name     string
	Decimals int
}{
	"0x41E94Eb019C0762f9Bfcf9Fb1E58725BfB0e7582": {Symbol: "USDC", Name: "USD Coin", Decimals: 6},
	"0x360ad4f9a9A8EFe9A8DCB5f461c4Cc1047E1Dcf9": {Symbol: "POL", Name: "Polygon", Decimals: 18},
	"0x0Fd9e8d3aF1aaee056EB9e802c3A762a667b1904": {Symbol: "LINK", Name: "Chainlink", Decimals: 18},
}

// getTokenBalancesPublicRPC handles token balance fetching for public RPC endpoints
func (c *AlchemyClient) getTokenBalancesPublicRPC(ctx context.Context, address string, chainID int, baseURL string) ([]*models.Balance, error) {
	var balances []*models.Balance
	
	// 1. Get native token balance using standard eth_getBalance
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
		return nil, fmt.Errorf("RPC error: %s", balanceResp.Error.Message)
	}

	// Add native balance if non-zero
	if balanceResp.Result != "0x0" && balanceResp.Result != "" {
		// Parse hex balance to big.Int
		balanceInt := new(big.Int)
		balanceInt.SetString(balanceResp.Result[2:], 16) // Remove 0x prefix
		
		balance := &models.Balance{
			ID:       uuid.New(),
			WalletID: uuid.New(),
			TokenID:  uuid.New(),
			Balance:  balanceInt.String(), // Store as decimal string
			Token: &models.Token{
				ID:       uuid.New(),
				Address:  "0x0000000000000000000000000000000000000000", // Native token
				ChainID:  chainID,
				Symbol:   "MATIC",
				Name:     "Polygon",
				Decimals: 18,
			},
		}
		balances = append(balances, balance)
	}

	// 2. Get ERC20 token balances for known tokens on Polygon Amoy
	if chainID == 80002 {
		for tokenAddr, tokenInfo := range polygonAmoyTokens {
			balance, err := c.getERC20Balance(ctx, address, tokenAddr, baseURL)
			if err != nil {
				logger.Error("Failed to get ERC20 balance", "token", tokenAddr, "error", err)
				continue
			}
			
			if balance != "0" {
				tokenBalance := &models.Balance{
					ID:       uuid.New(),
					WalletID: uuid.New(),
					TokenID:  uuid.New(),
					Balance:  balance,
					Token: &models.Token{
						ID:       uuid.New(),
						Address:  tokenAddr,
						ChainID:  chainID,
						Symbol:   tokenInfo.Symbol,
						Name:     tokenInfo.Name,
						Decimals: tokenInfo.Decimals,
					},
				}
				balances = append(balances, tokenBalance)
			}
		}
	}

	return balances, nil
}

// getERC20Balance fetches balance for a specific ERC20 token
func (c *AlchemyClient) getERC20Balance(ctx context.Context, walletAddress, tokenAddress, baseURL string) (string, error) {
	// ERC20 balanceOf method signature: 0x70a08231
	// Pad address to 32 bytes
	paddedAddress := fmt.Sprintf("0x70a08231%064s", walletAddress[2:])
	
	reqBody := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "eth_call",
		"params": []interface{}{
			map[string]string{
				"to":   tokenAddress,
				"data": paddedAddress,
			},
			"latest",
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "0", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL, strings.NewReader(string(reqBytes)))
	if err != nil {
		return "0", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "0", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var callResp struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&callResp); err != nil {
		return "0", fmt.Errorf("failed to decode response: %w", err)
	}

	if callResp.Error != nil {
		return "0", fmt.Errorf("RPC error: %s", callResp.Error.Message)
	}

	// Convert hex to decimal
	if callResp.Result == "0x" || callResp.Result == "" {
		return "0", nil
	}

	balance := new(big.Int)
	balance.SetString(callResp.Result[2:], 16) // Remove 0x prefix
	
	return balance.String(), nil
}

// getTransactionsPublicRPC handles transaction fetching for public RPC endpoints
func (c *AlchemyClient) getTransactionsPublicRPC(ctx context.Context, address string, chainID int, baseURL string) ([]*models.Transaction, error) {
	// For Polygon Amoy, we'll use Polygonscan API since standard RPC doesn't provide transaction history
	if chainID == 80002 {
		return c.getTransactionsFromPolygonscan(ctx, address, chainID)
	}
	
	// For other chains with public RPC, return empty for now
	// This would need to be implemented using block scanning or other methods
	return []*models.Transaction{}, nil
}

// getTransactionsFromPolygonscan fetches transactions from Polygonscan API for Polygon Amoy
func (c *AlchemyClient) getTransactionsFromPolygonscan(ctx context.Context, address string, chainID int) ([]*models.Transaction, error) {
	// Use Polygonscan API for Polygon Amoy testnet
	apiURL := fmt.Sprintf("https://api-amoy.polygonscan.com/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=50&sort=desc&apikey=YourApiKeyToken", address)
	
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var polygonscanResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			BlockNumber       string `json:"blockNumber"`
			TimeStamp         string `json:"timeStamp"`
			Hash              string `json:"hash"`
			From              string `json:"from"`
			To                string `json:"to"`
			Value             string `json:"value"`
			Gas               string `json:"gas"`
			GasPrice          string `json:"gasPrice"`
			IsError           string `json:"isError"`
			TxReceiptStatus   string `json:"txreceipt_status"`
			Input             string `json:"input"`
			ContractAddress   string `json:"contractAddress"`
			CumulativeGasUsed string `json:"cumulativeGasUsed"`
			GasUsed           string `json:"gasUsed"`
			Confirmations     string `json:"confirmations"`
			MethodId          string `json:"methodId"`
			FunctionName      string `json:"functionName"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&polygonscanResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if polygonscanResp.Status != "1" {
		logger.Error("Polygonscan API error", "message", polygonscanResp.Message)
		// Return empty slice instead of error to avoid breaking the app
		return []*models.Transaction{}, nil
	}

	// Convert to models.Transaction
	var transactions []*models.Transaction
	for _, tx := range polygonscanResp.Result {
		// Parse timestamp (Unix timestamp)
		timestamp, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
		blockNumber, _ := strconv.ParseInt(tx.BlockNumber, 10, 64)
		
		// Determine transaction type based on method and direction
		txType := "send"
		if strings.EqualFold(tx.From, address) {
			// Outgoing transaction
			if tx.FunctionName != "" {
				functionLower := strings.ToLower(tx.FunctionName)
				if strings.Contains(functionLower, "approve") {
					txType = "approve"
				} else if strings.Contains(functionLower, "stake") {
					txType = "send" // Staking is considered sending
				} else if strings.Contains(functionLower, "unstake") {
					txType = "receive" // Unstaking is considered receiving
				} else if strings.Contains(functionLower, "swap") {
					txType = "swap"
				} else {
					txType = "send"
				}
			} else {
				txType = "send"
			}
		} else {
			// Incoming transaction
			txType = "receive"
		}
		
		// Determine status
		status := "success"
		if tx.IsError == "1" || tx.TxReceiptStatus == "0" {
			status = "failed"
		}

		gasUsed, _ := strconv.ParseInt(tx.GasUsed, 10, 64)
		
		transaction := &models.Transaction{
			ID:          uuid.New(),
			Hash:        tx.Hash,
			ChainID:     chainID,
			FromAddress: tx.From,
			ToAddress:   &tx.To,
			Value:       &tx.Value,
			GasUsed:     &gasUsed,
			GasPrice:    &tx.GasPrice,
			BlockNumber: &blockNumber,
			Timestamp:   time.Unix(timestamp, 0),
			Status:      status,
			Type:        txType,
			Metadata: map[string]interface{}{
				"gas":             tx.Gas,
				"gasUsed":         tx.GasUsed,
				"methodId":        tx.MethodId,
				"functionName":    tx.FunctionName,
				"contractAddress": tx.ContractAddress,
			},
		}

		transactions = append(transactions, transaction)
	}

	logger.Info("Successfully fetched transactions from Polygonscan", 
		"address", address, 
		"transactionCount", len(transactions))

	return transactions, nil
}