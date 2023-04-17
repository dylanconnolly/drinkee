package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	user    = "dconnolly"
	dbname  = "drinkee"
	sslmode = "disable" // or verify-full
)

func CreatePostgresConnection() (*sqlx.DB, error) {
	connUrl := fmt.Sprintf("postgres://localhost:5432/%s?sslmode=%s", os.Getenv("POSTGRES_DBNAME"), os.Getenv("POSTGRES_SSLMODE"))
	// connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", user, dbname, sslmode)
	db, err := sqlx.Open("postgres", connUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

func SetLimitOffset(limit int, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	} else if offset > 0 {
		return fmt.Sprintf("OFFSET %d", offset)
	}
	return ""
}
