package main

import (
	"log"
	"net/http"

	"github.com/spikycham/finance/business"
	"github.com/spikycham/finance/network"
)

const PORT = "3000"

func main() {
	db, err := business.Connect("data.db")
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	r := business.NewRepository(db)
	s := business.NewService(r)
	h := business.NewHandler(s)
	log.Println(h)

	http.HandleFunc("/", welcome)

	log.Printf("Service is running at %s...", PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Printf("Failed to start serving: %v", err)
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	t := "Welcome to my personal financial records service. This is private service. A valid API key is requried to access endpoints."
	network.ResponseMessage(w, http.StatusOK, t)
}
