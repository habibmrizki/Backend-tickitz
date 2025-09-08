package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDb() (*pgxpool.Pool, error) {
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	connstring := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	return pgxpool.New(context.Background(), connstring)
}

func TesbDB(db *pgxpool.Pool) error {
	return db.Ping(context.Background())
}
