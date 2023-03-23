package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Drink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"desc"`
	Instructions string `json:"instructions"`
}

type NewDrinkRequest struct {
	Name             string            `json:"name" binding:"required"`
	Description      string            `json:"description" binding:"required"`
	Instructions     string            `json:"instructions" binding:"required"`
	DrinkIngredients []DrinkIngredient `json:"drink_ingredients" binding:"required"`
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

func (br *BaseRouter) createDrink(c *gin.Context) {
	var dr NewDrinkRequest

	err := c.ShouldBindJSON(&dr)
	if err != nil {
		c.String(http.StatusBadRequest, "can't bind: %s", err)
		return
	}

	// stmt, err := br.db.PrepareNamed("INSERT INTO drinks (name, description, instructions) VALUES (:name, :description, :instructions) RETURNING id")
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, "insert failed: %s", err)
	// 	return
	// }
	// err = stmt.Get(&drinkID, dr)

	// for _, id := range dr.IngredientIDs {
	// 	_, err := br.db.Exec("INSERT INTO drink_ingredients (drink_id, ingredient_id, measurement) VALUES ($1, $2, $3)", drinkID, id, "100 shots")
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, "drink ingredients insert failed: %s", err)
	// 		return
	// 	}
	// }

	fmt.Printf("%+v", dr.DrinkIngredients)

	br.db.MustExec(`
		WITH drink AS (
			INSERT INTO drinks (name, description, instructions)
			VALUES ($1, $2, $3)
			RETURNING id
		),
		ingredient_ids AS (
			SELECT id, name FROM ingredients WHERE name IN ('Tequila', 'Triple sec', 'Lime juice', 'Ice')
		),
		INSERT INTO drink_ingredients (drink_id, ingredient_id, measurement)
		SELECT drink.id, ingredient_ids.id, drink_ingredients.measurement
		FROM drink, ingredient_ids
		JOIN UNNEST($4::drink_ingredients[]) AS drink_ingredients ON drink_ingredients.name = ingredient_ids.name`, dr.Name, dr.Description, dr.Instructions, pq.Array(dr.DrinkIngredients))

	c.JSON(202, "added new drink")
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

	router.GET("/drinks", func(c *gin.Context) {
		br.getDrinks(c)
	})
	router.GET("/drinks/:id", func(c *gin.Context) {
		br.getDrinkByID(c)
	})
	router.POST("/drinks", func(c *gin.Context) {
		br.createDrink(c)
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
