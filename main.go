package main

import (
	"fmt"
	"log"

	"github.com/dylanconnolly/drinkee/postgres"
	"github.com/dylanconnolly/drinkee/router"
)

func main() {
	fmt.Println("connecting to postgres")
	db, err := postgres.CreatePostgresConnection()
	if err != nil {
		log.Fatal(err)
		fmt.Println("error connecting to postgres")
		return
	}
	fmt.Println("connection to postgres successful!")

	r := router.CreateNewRouter(db)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
