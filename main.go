package main

import (
	"net/http"
	"os"

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
		middleware.Auth(apiKey),
	)

	mux.HandleFunc("GET /", welcome)
	mux.HandleFunc("POST /api/v1/create", h.CreateItem)
	mux.HandleFunc("GET  /api/v1/items", h.GetItems)

	logger.Info("Service is running at " + PORT)
	if err := http.ListenAndServe(":"+PORT, handler); err != nil {
		logger.Error("Failed to start serving", err)
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	t := "Welcome to my personal financial records service. This is private service. A valid API key is requried to access endpoints."
	network.ResponseMessage(w, http.StatusOK, t)
}
