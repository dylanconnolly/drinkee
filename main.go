package main

import (
	"fmt"

	"github.com/dylanconnolly/drinkee/mongo"
	"github.com/dylanconnolly/drinkee/router"
)

func main() {
	fmt.Println("connecting to mongo...")
	err := mongo.CreateMongoConnection()
	if err != nil {
		fmt.Printf("error connecting to mongo: %s", err)
	}
	r := router.CreateRouter()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
