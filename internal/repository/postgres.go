package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func NewPostgresRepository() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/akademi?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("erro ao abrir conexão com o banco:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("erro ao conectar com o banco:", err)
	}

	log.Println("conectado ao Postgres com sucesso")

	return db
}
