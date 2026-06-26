package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spikycham/finance/business"
	"github.com/spikycham/finance/middleware"
	"github.com/spikycham/finance/network"
)

const PORT = "3000"

func main() {
	db, err := business.Connect("data.db")
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	r := business.NewSQLiteRepository(db)
	s := business.NewService(r)
	h := business.NewHandler(s)

	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load dotenv: %v", err)
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Print("no api key provided")
	}

	mux := http.NewServeMux()
	handler := middleware.Chain(
		mux,
		middleware.Auth(apiKey),
	)

	mux.HandleFunc("GET /", welcome)
	mux.HandleFunc("POST /api/v1/create", h.CreateItem)
	mux.HandleFunc("GET  /api/v1/items", h.GetItems)

	log.Printf("Service is running at %s...", PORT)
	if err := http.ListenAndServe(":"+PORT, handler); err != nil {
		log.Printf("Failed to start serving: %v", err)
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	t := "Welcome to my personal financial records service. This is private service. A valid API key is requried to access endpoints."
	network.ResponseMessage(w, http.StatusOK, t)
}
