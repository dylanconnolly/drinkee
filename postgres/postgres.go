package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	user    = "dconnolly"
	dbname  = "drinkee"
	sslmode = "disable" // or verify-full
)

func CreatePostgresConnection() (*sqlx.DB, error) {
	connUrl := fmt.Sprintf("postgres://localhost:5432/drinkee?sslmode=%s", sslmode)
	// connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", user, dbname, sslmode)
	db, err := sqlx.Open("postgres", connUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
