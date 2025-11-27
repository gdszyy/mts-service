package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdsZyy/mts-service/internal/api"
	"github.com/gdsZyy/mts-service/internal/client"
	"github.com/gdsZyy/mts-service/internal/config"
	"github.com/gdsZyy/mts-service/internal/service"
)

func main() {
	log.Println("Starting MTS Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Configuration loaded: Production=%v, Port=%s", cfg.Production, cfg.Port)

	// Auto-fetch Bookmaker ID if not provided and AccessToken is available
	if cfg.BookmakerID == "" && cfg.AccessToken != "" {
		log.Println("Bookmaker ID not provided, attempting to fetch from whoami.xml...")
		bookmakerID, _, err := client.FetchBookmakerInfo(cfg.AccessToken, cfg.UOFAPIBaseURL)
		if err != nil {
			log.Printf("Warning: Failed to fetch Bookmaker Info from whoami.xml: %v. Proceeding without auto-configuration.", err)
		} else {
			cfg.BookmakerID = bookmakerID
			log.Printf("Bookmaker ID fetched successfully: %s", bookmakerID)
		}
	}

	// Check if BookmakerID is still empty after auto-fetch attempt
	if cfg.BookmakerID == "" {
		log.Println("Bookmaker ID is still empty. MTS service will not start.")
		return
	}

	// Create MTS service
	mtsService := service.NewMTSService(cfg)

	// Start MTS service
	if err := mtsService.Start(); err != nil {
		log.Fatalf("Failed to start MTS service: %v", err)
	}

	// Create API handler
	handler := api.NewHandler(mtsService, cfg)

	// Setup routes
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("/health", handler.HealthCheck)
	
	// Legacy endpoint (for backward compatibility)
	mux.HandleFunc("/api/tickets", handler.PlaceTicket)
	
	// New bet type endpoints
	mux.HandleFunc("/api/bets/single", handler.PlaceSingleBet)
	mux.HandleFunc("/api/bets/accumulator", handler.PlaceAccumulatorBet)
	mux.HandleFunc("/api/bets/system", handler.PlaceSystemBet)
	mux.HandleFunc("/api/bets/banker-system", handler.PlaceBankerSystemBet)
	mux.HandleFunc("/api/bets/preset", handler.PlacePresetSystemBet)
	mux.HandleFunc("/api/bets/multi", handler.PlaceMultiBet)
	
	// Cashout endpoint
	mux.HandleFunc("/api/cashout", handler.RequestCashout)
	
	// Root endpoint with API documentation
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"service": "mts-service",
			"version": "2.0.0",
			"endpoints": {
				"health": "/health",
				"legacy": "/api/tickets",
				"bets": {
					"single": "/api/bets/single",
					"accumulator": "/api/bets/accumulator",
					"system": "/api/bets/system",
					"banker_system": "/api/bets/banker-system",
					"preset": "/api/bets/preset",
					"multi": "/api/bets/multi"
				},
				"cashout": "/api/cashout"
			}
		}`))
	})

	// Start HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: enableCORS(mux),
	}

	go func() {
		log.Printf("HTTP server listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down gracefully...")
	mtsService.Stop()
	server.Close()
	log.Println("Service stopped")
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

