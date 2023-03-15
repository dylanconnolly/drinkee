package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	user   = "dconnolly"
	dbname = "drinkee"
)

func CreatePostgresConnection() error {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=verify-full", user, dbname)
	_, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
