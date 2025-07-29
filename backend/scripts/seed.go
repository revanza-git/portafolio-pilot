package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Seed data
	if err := seedData(ctx, pool); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("Database seeded successfully!")
}

func seedData(ctx context.Context, pool *pgxpool.Pool) error {
	// Start transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Seed users
	userIDs := make([]uuid.UUID, 0)
	addresses := []string{
		"0x1234567890123456789012345678901234567890",
		"0x0987654321098765432109876543210987654321",
		"0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
	}

	for _, addr := range addresses {
		var userID uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO users (address, nonce) 
			VALUES ($1, $2) 
			ON CONFLICT (address) DO UPDATE SET updated_at = NOW()
			RETURNING id
		`, addr, generateNonce()).Scan(&userID)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
		userIDs = append(userIDs, userID)
		log.Printf("Created user: %s", addr)
	}

	// Seed wallets
	walletIDs := make([]uuid.UUID, 0)
	for i, userID := range userIDs {
		// Primary wallet on Ethereum
		var walletID uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO wallets (user_id, address, chain_id, label, is_primary)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (address, chain_id) DO UPDATE SET updated_at = NOW()
			RETURNING id
		`, userID, addresses[i], 1, "Main Wallet", true).Scan(&walletID)
		if err != nil {
			return fmt.Errorf("failed to insert wallet: %w", err)
		}
		walletIDs = append(walletIDs, walletID)

		// Secondary wallet on Polygon
		err = tx.QueryRow(ctx, `
			INSERT INTO wallets (user_id, address, chain_id, label, is_primary)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (address, chain_id) DO UPDATE SET updated_at = NOW()
			RETURNING id
		`, userID, addresses[i], 137, "Polygon Wallet", false).Scan(&walletID)
		if err != nil {
			return fmt.Errorf("failed to insert wallet: %w", err)
		}
	}

	// Seed tokens
	tokens := []struct {
		address  string
		chainID  int
		symbol   string
		name     string
		decimals int
		price    float64
	}{
		{"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", 1, "WETH", "Wrapped Ether", 18, 2345.67},
		{"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", 1, "USDC", "USD Coin", 6, 1.0},
		{"0x6B175474E89094C44Da98b954EedeAC495271d0F", 1, "DAI", "Dai Stablecoin", 18, 1.0},
		{"0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599", 1, "WBTC", "Wrapped Bitcoin", 8, 45678.90},
		{"0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619", 137, "WETH", "Wrapped Ether", 18, 2345.67},
		{"0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", 137, "USDC", "USD Coin", 6, 1.0},
	}

	tokenIDs := make(map[string]uuid.UUID)
	for _, token := range tokens {
		var tokenID uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO tokens (address, chain_id, symbol, name, decimals, price_usd, price_change_24h, market_cap, last_updated)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (address, chain_id) DO UPDATE SET 
				price_usd = $6,
				price_change_24h = $7,
				market_cap = $8,
				last_updated = $9,
				updated_at = NOW()
			RETURNING id
		`, token.address, token.chainID, token.symbol, token.name, token.decimals, 
		   token.price, 2.5, token.price*1000000, time.Now()).Scan(&tokenID)
		if err != nil {
			return fmt.Errorf("failed to insert token: %w", err)
		}
		tokenIDs[fmt.Sprintf("%s-%d", token.address, token.chainID)] = tokenID
		log.Printf("Created token: %s on chain %d", token.symbol, token.chainID)
	}

	// Seed some balances
	for i, walletID := range walletIDs[:3] { // First 3 wallets (Ethereum mainnet)
		// WETH balance
		wethKey := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2-1"
		if tokenID, ok := tokenIDs[wethKey]; ok {
			balance := fmt.Sprintf("%d", (i+1)*1000000000000000000) // 1-3 ETH
			balanceUSD := float64(i+1) * 2345.67
			_, err := tx.Exec(ctx, `
				INSERT INTO balances (wallet_id, token_id, balance, balance_usd)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (wallet_id, token_id) DO UPDATE SET
					balance = $3,
					balance_usd = $4,
					updated_at = NOW()
			`, walletID, tokenID, balance, balanceUSD)
			if err != nil {
				return fmt.Errorf("failed to insert balance: %w", err)
			}
		}

		// USDC balance
		usdcKey := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48-1"
		if tokenID, ok := tokenIDs[usdcKey]; ok {
			balance := fmt.Sprintf("%d", (i+1)*5000000000) // 5000-15000 USDC
			balanceUSD := float64((i + 1) * 5000)
			_, err := tx.Exec(ctx, `
				INSERT INTO balances (wallet_id, token_id, balance, balance_usd)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (wallet_id, token_id) DO UPDATE SET
					balance = $3,
					balance_usd = $4,
					updated_at = NOW()
			`, walletID, tokenID, balance, balanceUSD)
			if err != nil {
				return fmt.Errorf("failed to insert balance: %w", err)
			}
		}
	}

	// Seed some transactions
	for i, addr := range addresses {
		hash := fmt.Sprintf("0x%064d", i+1)
		_, err := tx.Exec(ctx, `
			INSERT INTO transactions (
				hash, chain_id, from_address, to_address, value,
				gas_used, gas_price, gas_fee_usd, block_number,
				timestamp, status, type
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT (hash) DO NOTHING
		`, hash, 1, addr, addresses[(i+1)%3], "1000000000000000000",
		   21000, "20000000000", 2.50, 18500000+i,
		   time.Now().Add(-time.Duration(i)*24*time.Hour), "confirmed", "send")
		if err != nil {
			return fmt.Errorf("failed to insert transaction: %w", err)
		}
		log.Printf("Created transaction: %s", hash)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func generateNonce() string {
	return fmt.Sprintf("nonce-%d", time.Now().UnixNano())
}