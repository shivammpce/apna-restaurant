package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectToDB() *sql.DB {
	userName := os.Getenv("DB_USER")
	userPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbURI := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", userName, userPassword, dbName)
	dbConfig, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("error loading database config", err)
	}
	return dbConfig
}
