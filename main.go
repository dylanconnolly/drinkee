package main

import (
	"fmt"
	"log"

	"github.com/dylanconnolly/drinkee/http"
	"github.com/dylanconnolly/drinkee/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const DefaultConfigPath = "~/.env"

type Main struct {
	DB         *sqlx.DB
	HTTPServer *http.Server
}

func CreateMain() *Main {
	db, _ := postgres.CreatePostgresConnection()
	return &Main{
		DB:         db,
		HTTPServer: http.NewServer(),
	}
}

func main() {
	// fmt.Println("db username env: ", os.Getenv("POSTGRES_USERNAME"))
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}

	fmt.Println("creating main")
	m := CreateMain()
	drinkService := postgres.NewDrinkService(m.DB)
	m.HTTPServer.DrinkService = drinkService
	m.HTTPServer.Serve()
}
