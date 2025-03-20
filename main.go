package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
)

var (
	dsn              = flag.String("dsn", "", "PostgreSQL connection string")
	webListenAddress = flag.String("web.listen-address", "", "Address to listen on for health checks")
)

func main() {
	flag.Parse()

	if *dsn == "" || *webListenAddress == "" {
		log.Fatal("Both --dsn and --web.listen-address must be specified")
	}

	http.HandleFunc("/health", healthHandler)
	log.Printf("Starting server on %s", *webListenAddress)
	log.Fatal(http.ListenAndServe(*webListenAddress, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt database connection
	conn, err := pgx.Connect(ctx, *dsn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database connection failed: %v", err), http.StatusServiceUnavailable)
		return
	}
	defer conn.Close(ctx)

	// Execute health check query
	var result int
	err = conn.QueryRow(ctx, "SELECT 1;").Scan(&result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Health query failed: %v", err), http.StatusServiceUnavailable)
		return
	}

	if result != 1 {
		http.Error(w, "Unexpected result from health query", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}
