package http

import (
	"net/http"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleGetDrinks(c *gin.Context) {
	drinks, err := s.DrinkService.FindDrinks(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "error getting drinks: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, drinks)
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
