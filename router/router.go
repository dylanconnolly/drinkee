package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/dylanconnolly/drinkee/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Drink struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName" db:"display_name"`
	Description  string `json:"description,omitempty"`
	Instructions string `json:"instructions"`
}

type DrinkResponse struct {
	Drink
	DrinkIngredients DrinkIngredientSlice `json:"drinkIngredients" db:"drink_ingredients"`
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
	DisplayName string `json:"displayName" db:"display_name"`
}

type DrinkIngredient struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName" db:"display_name"`
	Measurement string `json:"measurement"`
}

type DrinkIngredientSlice []DrinkIngredient

func (dis *DrinkIngredientSlice) Scan(src interface{}) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return nil
	}
	return json.Unmarshal(data, dis)
}

// func (dis DrinkIngredientSlice) Value() (driver.Value, error) {
// 	return json.Marshal(dis)
// }

type BaseRouter struct {
	db           *sqlx.DB
	DrinkService *postgres.DrinkService
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

// func (br *BaseRouter) getDrinks(c *gin.Context) {
// 	var drinks []DrinkResponse

// 	queryStr := `
// 	SELECT d.id, d.name, d.display_name, d.description, d.instructions, json_agg(json_build_object('name', i.name, 'displayName', i.display_name, 'measurement', di.measurement)) as drink_ingredients
// 	FROM drinks d
// 	JOIN drink_ingredients di ON di.drink_id=d.id
// 	JOIN ingredients i ON di.ingredient_id=i.id
// 	GROUP BY d.id, d.name ORDER BY d.name;
// 	`

// 	err := br.db.Select(&drinks, queryStr)

// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "err= %s", err)
// 	}

// 	c.IndentedJSON(http.StatusOK, drinks)
// }

func (br *BaseRouter) getDrinks(c *gin.Context) {
	f := drinkee.DrinkFilter{}
	drinks, err := br.DrinkService.FindDrinks(c, f)
	if err != nil {
		c.String(http.StatusInternalServerError, "error getting drinks: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, drinks)
	fmt.Println("youre in OLD func")
}

func (br *BaseRouter) getDrinkByID(c *gin.Context) {
	id := c.Param("id")
	var drink Drink
	var drinkIngredients []DrinkIngredient

	err := br.db.Get(&drink, "SELECT id, name, display_name, description, instructions FROM drinks WHERE id=$1", id)
	if err != nil {
		c.String(http.StatusBadRequest, "No drink with that id:", err)
		return
	}

	queryStr :=
		`
		SELECT ingredients.name, ingredients.display_name, drink_ingredients.measurement
		FROM drink_ingredients JOIN ingredients ON ingredients.id = drink_ingredients.ingredient_id
		WHERE drink_id=$1
		`

	err = br.db.Select(&drinkIngredients, queryStr, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "couldnt find drink ingredients idk", err)
	}

	drinkData := DrinkResponse{drink, drinkIngredients}

	c.IndentedJSON(http.StatusOK, drinkData)
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
	SELECT ingredients.name, ingredients.display_name, drink_ingredients.measurement
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

	err := br.db.Select(&ingredients, "SELECT id, name, display_name FROM ingredients ORDER BY name")
	if err != nil {
		c.String(http.StatusInternalServerError, "err getting ingredients: %s", err)
	}

	c.IndentedJSON(http.StatusOK, ingredients)
}

func (br *BaseRouter) getIngredientByID(c *gin.Context) {
	id := c.Param("id")
	var ingredient Ingredient

	err := br.db.Get(&ingredient, "SELECT id, name, display_name FROM ingredients WHERE id=$1", id)
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
	// var drink Drink
	var drinks []DrinkResponse

	err := c.ShouldBindJSON(&ingredientList)
	if err != nil {
		c.String(http.StatusBadRequest, "couldn't bind to list of ingredients", err)
		return
	}

	for _, ingredient := range ingredientList.Ingredients {
		ingredientIDs = append(ingredientIDs, ingredient.ID)
	}

	// queryStr := `
	// 	SELECT id,name,description,instructions
	// 	FROM
	// 		(SELECT d.*, COUNT(*) AS ingredients_present,
	// 		(SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) AS total_ingredients
	// 		FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id = ANY($1) GROUP BY d.id) AS joiny
	// 	WHERE ingredients_present=total_ingredients`

	queryStr := `SELECT md.id,md.name,md.display_name,md.description,md.instructions, ij.drink_ingredients
		FROM 
			(SELECT d.*, COUNT(*) AS ingredients_present,
			(SELECT COUNT(*) FROM drink_ingredients WHERE drink_ingredients.drink_id=d.id) AS total_ingredients 
			FROM drinks d JOIN drink_ingredients di ON di.drink_id=d.id WHERE di.ingredient_id = ANY($1) GROUP BY d.id) AS md 
      JOIN (SELECT d.id, json_agg(json_build_object('name', i.name, 'displayName', i.display_name, 'measurement', di.measurement)) as drink_ingredients 
            FROM drinks d 
            JOIN drink_ingredients di ON di.drink_id=d.id
            JOIN ingredients i ON di.ingredient_id=i.id 
            GROUP BY d.id, d.name ) AS ij ON ij.id=md.id
		WHERE ingredients_present=total_ingredients
		ORDER BY md.name;`

	err = br.db.Select(&drinks, queryStr, pq.Array(ingredientIDs))
	if err != nil {
		c.String(http.StatusInternalServerError, "error generating cocktails: %s", err)
		return
	}
	// rows, err := br.db.Queryx(queryStr, pq.Array(ingredientIDs))

	// for rows.Next() {
	// 	err = rows.StructScan(&drink)
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, "error scanning drink into struct: ", err)
	// 		return
	// 	}
	// 	drinks = append(drinks, drink)
	// }

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
		db:           db,
		DrinkService: postgres.NewDrinkService(db),
	}

	router := gin.Default()

	router.Use(cors.Default())

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"POST"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

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
