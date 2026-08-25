package main

import (
	"log"
	"net/http"

	"github.com/TakedaB/akademi-api/cmd/internal/repository"
)

func main() {
	db := repository.NewPostgresRepository()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)

	log.Println("servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
