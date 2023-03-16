package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Drink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"desc"`
	Instructions string `json:"instructions"`
}

type BaseRouter struct {
	db *sqlx.DB
}

func (br *BaseRouter) getDrinks(c *gin.Context) {
	// drinks := []Drink{}
	var drinks []Drink

	br.db.Select(&drinks, "SELECT * FROM drinks")

	c.IndentedJSON(http.StatusOK, drinks)
}

func (br *BaseRouter) getDrinkByID(c *gin.Context) {
	id := c.Param("id")
	var drink Drink

	err := br.db.Get(&drink, "SELECT * FROM drinks WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusOK, "No drink with that id")
		return
	}

	c.IndentedJSON(http.StatusOK, drink)
}

func CreateNewRouter(db *sqlx.DB) *gin.Engine {
	br := &BaseRouter{
		db: db,
	}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/drinks", func(c *gin.Context) {
		br.getDrinks(c)
	})
	router.GET("/drinks/:id", func(c *gin.Context) {
		br.getDrinkByID(c)
	})

	return router
}
