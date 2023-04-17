package http

import (
	"net/http"
	"strconv"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/gin-gonic/gin"
)

type IngredientListRequest struct {
	Ingredients []drinkee.Ingredient `json:"ingredients"`
}

func (s *Server) handleGetDrinks(c *gin.Context) {
	drinks, err := s.DrinkService.FindDrinks(c)
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
	c.String(http.StatusAccepted, "non strict generate cocktails")
}

func (s *Server) handleGetIngredients(c *gin.Context) {
	ingredients, err := s.DrinkService.FindIngredients(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "error getting ingredients: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, ingredients)
}
