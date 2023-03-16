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
	return r
}
