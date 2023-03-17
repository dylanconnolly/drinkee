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

type Ingredient struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DrinkIngredient struct {
	Name        string `json:"name"`
	Measurement string `json:"measurement"`
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

func (br *BaseRouter) getDrinkIngredients(c *gin.Context) {
	drinkID := c.Param("id")
	var drinkIngredients []DrinkIngredient

	queryStr := `
	SELECT ingredients.name, drink_ingredients.measurement
	FROM drink_ingredients JOIN ingredients ON ingredients.id = drink_ingredients.ingredient_id
	WHERE drink_id=$1`

	err := br.db.Select(&drinkIngredients, queryStr, drinkID)
	if err != nil {
		c.JSON(http.StatusOK, "No drink ingredients for that drink")
		return
	}

	c.IndentedJSON(http.StatusOK, drinkIngredients)
}

func (br *BaseRouter) getIngredients(c *gin.Context) {
	var ingredients []Ingredient

	br.db.Select(&ingredients, "SELECT * FROM ingredients")

	c.IndentedJSON(http.StatusOK, ingredients)
}

func (br *BaseRouter) getIngredientByID(c *gin.Context) {
	id := c.Param("id")
	var ingredient Ingredient

	err := br.db.Get(&ingredient, "SELECT * FROM ingredients WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusOK, "No drink with that id")
		return
	}

	c.IndentedJSON(http.StatusOK, ingredient)
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
	router.GET("drinks/:id/ingredients", func(c *gin.Context) {
		br.getDrinkIngredients(c)
	})
	router.GET("/ingredients", func(c *gin.Context) {
		br.getIngredients(c)
	})
	router.GET("/ingredients/:id", func(c *gin.Context) {
		br.getIngredientByID(c)
	})

	return router
}
