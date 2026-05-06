package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {

	log.Printf("Iniciando conexão com Banco de Dados...")

	psqlInfo := getPsqlInfo()

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("ERRO: Variável de ambiente %s não definida!", key)
	}
	return value
}

func getPsqlInfo() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnvOrPanic("POSTGRES_HOST"),
		getEnvOrPanic("POSTGRES_PORT"),
		getEnvOrPanic("POSTGRES_USER"),
		getEnvOrPanic("POSTGRES_PASSWORD"),
		getEnvOrPanic("POSTGRES_DB"),
	)
}
