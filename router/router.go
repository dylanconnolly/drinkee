package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Drink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
}

type CreateDrinkRequest struct {
	Name             string            `json:"name" binding:"required"`
	DisplayName      string            `json:"displayName" binding:"required"`
	Description      string            `json:"description"`
	Instructions     string            `json:"instructions" binding:"required"`
	DrinkIngredients []DrinkIngredient `json:"drinkIngredients" binding:"required"`
}

type Ingredient struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type DrinkIngredient struct {
	Name        string `json:"name"`
	Measurement string `json:"measurement"`
}

type BaseRouter struct {
	db *sqlx.DB
}

type IngredientsListRequest struct {
	Ingredients []Ingredient `json:"ingredients"`
}

type CreateIngredientRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"displayName" binding:"required" db:"display_name"`
}

type CreateIngredientsListRequest struct {
	Ingredients []CreateIngredientRequest `json:"ingredients" binding:"required"`
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
	var dr CreateDrinkRequest
	var ingredientNames []string

	err := c.ShouldBindJSON(&dr)
	if err != nil {
		c.String(http.StatusBadRequest, "can't bind: %s", err)
		return
	}

	for _, di := range dr.DrinkIngredients {
		ingredientNames = append(ingredientNames, di.Name)
	}

	drinkIngredientsJSON, err := json.Marshal(dr.DrinkIngredients)

	_, err = br.db.Exec(`
	WITH drink AS (
		INSERT INTO drinks (name, display_name, description, instructions)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	),
	ingredient_ids AS (
		SELECT id, name FROM ingredients WHERE name = ANY($5)
	),
	ingredient_data AS (
		SELECT * FROM json_populate_recordset(null::ingredient_data, $6)
	)
	INSERT INTO drink_ingredients (drink_id, ingredient_id, measurement)
	SELECT drink.id, ingredient_ids.id, ingredient_data.measurement
	FROM drink, ingredient_ids, ingredient_data
	WHERE ingredient_ids.name = ingredient_data.name`, dr.Name, dr.DisplayName, dr.Description, dr.Instructions, pq.Array(ingredientNames), string(drinkIngredientsJSON))

	if err != nil {
		c.String(http.StatusInternalServerError, "error adding drink: %s", err)
		return
	}

	c.String(202, "added new drink\n")
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

func (br *BaseRouter) createIngredientsFromList(c *gin.Context) {
	var ilr CreateIngredientsListRequest

	err := c.ShouldBindJSON(&ilr)
	if err != nil {
		c.String(http.StatusBadRequest, "can't bind request to ingredient list: %s", err)
		fmt.Printf("can't bind request to ingredient list: %+v", err)
		return
	}

	_, err = br.db.NamedExec(`INSERT INTO ingredients (name, display_name) VALUES (:name, :display_name)`, ilr.Ingredients)
	if err != nil {
		c.String(http.StatusInternalServerError, "error adding ingredients: ", err)
		return
	}

	c.String(http.StatusAccepted, "added ingredients")
}

func (br *BaseRouter) generateCocktails(c *gin.Context) {
	c.String(http.StatusOK, "unstrict generate cocktails")
}

func (br *BaseRouter) generateCocktailsStrict(c *gin.Context) {
	var ingredientList IngredientsListRequest
	var ingredientIDs []int
	var drink Drink
	var drinks []Drink

	err := c.ShouldBindJSON(&ingredientList)
	if err != nil {
		c.String(http.StatusBadRequest, "couldn't bind to list of ingredients", err)
	}

	for _, ingredient := range ingredientList.Ingredients {
		ingredientIDs = append(ingredientIDs, ingredient.ID)
	}

	queryStr := `SELECT id,name,description,instructions FROM (SELECT d.*, COUNT(*) AS ingredients_present, (SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) AS total_ingredients FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id = ANY($1) GROUP BY d.id) AS joiny WHERE ingredients_present=total_ingredients`
	rows, err := br.db.Queryx(queryStr, pq.Array(ingredientIDs))
	for rows.Next() {
		err = rows.StructScan(&drink)
		if err != nil {
			c.String(http.StatusInternalServerError, "error scanning drink into struct: ", err)
		}
		drinks = append(drinks, drink)
	}

	c.IndentedJSON(http.StatusOK, drinks)
}

func (br *BaseRouter) createIngredientsList(c *gin.Context) {
	var ingredients IngredientsListRequest

	err := c.ShouldBindJSON(&ingredients)
	if err != nil {
		c.String(http.StatusBadRequest, "couldn't bind to list of ingredients", err)
	}

	fmt.Printf("ingredient list: %v", ingredients)
}

func CreateNewRouter(db *sqlx.DB) *gin.Engine {
	br := &BaseRouter{
		db: db,
	}

	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/drinks", func(c *gin.Context) {
				br.getDrinks(c)
			})
			v1.GET("/drinks/:id", func(c *gin.Context) {
				br.getDrinkByID(c)
			})
			v1.POST("/drinks", func(c *gin.Context) {
				br.createDrink(c)
			})
			v1.GET("drinks/:id/ingredients", func(c *gin.Context) {
				br.getDrinkIngredients(c)
			})
			v1.GET("/ingredients", func(c *gin.Context) {
				br.getIngredients(c)
			})
			v1.POST("/ingredients", func(c *gin.Context) {
				br.createIngredientsFromList(c)
			})
			v1.GET("/ingredients/:id", func(c *gin.Context) {
				br.getIngredientByID(c)
			})
			v1.POST("/generateCocktails", func(c *gin.Context) {
				strict := c.Query("strict")
				if strict == "true" {
					br.generateCocktailsStrict(c)
				} else {
					br.generateCocktails(c)
				}
			})
			v1.POST("/ingredientsList", func(c *gin.Context) {
				br.createIngredientsList(c)
			})
		}
	}

	return router
}
