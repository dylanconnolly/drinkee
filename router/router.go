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

func getDrinks(c *gin.Context, db *sqlx.DB) {
	// drinks := []Drink{}
	var drinks []Drink

	db.Select(&drinks, "SELECT * FROM drinks")

	c.IndentedJSON(http.StatusOK, drinks)
}

func getDrinkByID(c *gin.Context, db *sqlx.DB) {
	id := c.Param("id")
	var drink Drink

	err := db.Get(&drink, "SELECT * FROM drinks WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusOK, "No drink with that id")
		return
	}

	c.IndentedJSON(http.StatusOK, drink)
}

func CreateRouter(db *sqlx.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/drinks", func(c *gin.Context) {
		getDrinks(c, db)
	})
	r.GET("/drinks/:id", func(c *gin.Context) {
		getDrinkByID(c, db)
	})

	return r
}
