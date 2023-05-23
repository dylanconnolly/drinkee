package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/gin-gonic/gin"
)

type IngredientListRequest struct {
	Ingredients []drinkee.Ingredient `json:"ingredients"`
}

func (s *Server) handleGetDrinks(c *gin.Context) {
	f := buildFilter(c)

	fmt.Printf("filter binding: %+v", f)

	drinks, err := s.DrinkService.FindDrinks(c, f)
	if err != nil {
		c.String(http.StatusInternalServerError, "error getting drinks: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, drinks)
}

func (s *Server) handleGetDrinkByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID format")
		return
	}

	drink, err := s.DrinkService.FindDrinkByID(c, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching drink: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, drink)
}

func (s *Server) handleCreateDrink(c *gin.Context) {
	var createDrink drinkee.CreateDrink

	if err := c.ShouldBindJSON(&createDrink); err != nil {
		c.String(http.StatusBadRequest, "invalid JSON in request body: %s", err)
		return
	}

	err := s.DrinkService.CreateDrink(c, &createDrink)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating drink: %s", err)
		return
	}

	c.String(http.StatusAccepted, "added new drink! \n")
}

func (s *Server) handleGenerateDrinks(c *gin.Context) {
	var ingredientList IngredientListRequest
	err := c.ShouldBindJSON(&ingredientList)
	if err != nil {
		c.String(http.StatusBadRequest, "couldn't bind to list of ingredients", err)
		return
	}

	ingredients := ingredientList.Ingredients

	strict := c.Query("strict")
	if strict == "true" {
		drinks, err := s.DrinkService.GenerateDrinks(c, ingredients)
		if err != nil {
			c.String(http.StatusInternalServerError, "error generating drinks: %s", err)
			return
		}

		c.IndentedJSON(http.StatusAccepted, drinks)
		return
	}
	drinks, err := s.DrinkService.GenerateNonStrictDrinks(c, ingredients)
	if err != nil {
		c.String(http.StatusInternalServerError, "error generating non strict drinks: %s", err)
	}
	c.IndentedJSON(http.StatusAccepted, drinks)
}

func (s *Server) handleGetIngredients(c *gin.Context) {
	ingredients, err := s.DrinkService.FindIngredients(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "error getting ingredients: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, ingredients)
}

func buildFilter(c *gin.Context) drinkee.DrinkFilter {
	var f drinkee.DrinkFilter

	c.ShouldBindJSON(&f)
	// if err := c.ShouldBindJSON(&f); err != nil {
	// 	fmt.Printf("error decoding filter in request body: %s", err)
	// }

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 100
	}
	skip, err := strconv.Atoi(c.Query("skip"))
	if err != nil {
		skip = 0
	}

	filter := drinkee.DrinkFilter{
		Limit: limit,
		Skip:  skip,
		Name:  f.Name,
		ID:    f.ID,
	}

	return filter
}
