package http

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) GenerateRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/drinks/:id", func(c *gin.Context) {
				s.handleGetDrinkByID(c)
			})
			v1.GET("/drinks", func(c *gin.Context) {
				s.handleGetDrinks(c)
			})
			v1.POST("/drinks", func(c *gin.Context) {
				s.handleCreateDrink(c)
			})
			v1.POST("/generateDrinks", func(c *gin.Context) {
				s.handleGenerateDrinks(c)
			})
			v1.GET("/ingredients", func(c *gin.Context) {
				s.handleGetIngredients(c)
			})
		}
	}
}
