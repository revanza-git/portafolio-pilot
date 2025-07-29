package main

import (
	"context"
	"fmt"
	"log"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Test different connection approaches
	connStrings := []struct{
		name string
		url  string
	}{
		{"127.0.0.1", "postgresql://defi:defi123@127.0.0.1:5432/defi_dashboard?sslmode=disable"},
		{"Container IP", "postgresql://defi:defi123@172.17.0.2:5432/defi_dashboard?sslmode=disable"},
		{"Docker Desktop IP", "postgresql://defi:defi123@192.168.65.254:5432/defi_dashboard?sslmode=disable"},
	}

	for _, cs := range connStrings {
		fmt.Printf("\nTesting %s:\n", cs.name)
		
		conn, err := pgx.Connect(context.Background(), cs.url)
		if err != nil {
			log.Printf("  Failed: %v", err)
			continue
		}
		
		var result int
		err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&result)
		if err != nil {
			log.Printf("  Query failed: %v", err)
		} else {
			fmt.Printf("  Success! Connected via %s\n", cs.name)
		}
		
		conn.Close(context.Background())
	}
}