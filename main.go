// PERF: recognize the service is shutting down, and response with notification.
// PERF: update and delete items from records.

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/spikycham/finance/business"
	"github.com/spikycham/finance/logger"
	"github.com/spikycham/finance/middleware"
	"github.com/spikycham/finance/network"
)

const PORT = "3000"

func main() {
	db, err := business.Connect("data.db")
	if err != nil {
		logger.Error("Failed to connect to database", err)
	}

	r := business.NewSQLiteRepository(db)
	s := business.NewService(r)
	h := business.NewHandler(s)

	if err := godotenv.Load(); err != nil {
		logger.Error("failed to load dotenv", err)
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		logger.Warn("no api key provided")
	}

	mux := http.NewServeMux()
	handler := middleware.Chain(
		mux,
		middleware.Logger,
		middleware.CORS,
		middleware.Auth(apiKey),
	)

	mux.HandleFunc("GET /", welcome)
	mux.HandleFunc("POST /api/v1/create", h.CreateItem)
	mux.HandleFunc("GET  /api/v1/items", h.GetYearlyItems)

	srv := &http.Server{
		Addr:    ":" + PORT,
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("Failed to start serving", err)
		}
	}()

	logger.Info("Service is running at " + PORT)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Shutdown failed", err)
	}

	db.Close()
	logger.Info("Service stopped")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	t := "Welcome to my personal financial records service. This is private service. A valid API key is requried to access endpoints."
	network.ResponseMessage(w, http.StatusOK, t)
}
