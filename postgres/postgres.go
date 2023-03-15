package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	// user   = "notauser"
	user    = "dconnolly"
	dbname  = "drinkee"
	sslmode = "disable" // or verify-full
)

func CreatePostgresConnection() error {
	connUrl := fmt.Sprintf("postgres://localhost:5432/drinkee?sslmode=%s", sslmode)
	// connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", user, dbname, sslmode)
	db, err := sql.Open("postgres", connUrl)

	if err != nil {
		log.Fatal(err)
		return err
	}

	insertStr := `INSERT INTO drinks (name, description, instructions) VALUES ('test drink', 'test description', 'test instructions')`

	_, err = db.Exec(insertStr)
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	return nil
}
