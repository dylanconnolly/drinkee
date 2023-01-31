package main

import (
	"github.com/dylanconnolly/drinkee/router"
)

func main() {
	r := router.CreateRouter()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
