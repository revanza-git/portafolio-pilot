package main

import (
	"context"
	"fmt"
	"log"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Test different connection strings
	connStrings := []string{
		"postgresql://defi:defi123@localhost:5432/defi_dashboard?sslmode=disable",
		"postgresql://defi:defi123@127.0.0.1:5432/defi_dashboard?sslmode=disable",
		"postgres://defi:defi123@localhost:5432/defi_dashboard?sslmode=disable",
		"host=localhost port=5432 user=defi password=defi123 dbname=defi_dashboard sslmode=disable",
	}

	for i, connStr := range connStrings {
		fmt.Printf("\nTesting connection string %d:\n", i+1)
		fmt.Printf("  %s\n", connStr)
		
		conn, err := pgx.Connect(context.Background(), connStr)
		if err != nil {
			log.Printf("  Failed: %v", err)
			continue
		}
		
		var result int
		err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&result)
		if err != nil {
			log.Printf("  Query failed: %v", err)
		} else {
			fmt.Printf("  Success! Result: %d\n", result)
		}
		
		conn.Close(context.Background())
	}
}