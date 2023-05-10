package http

import (
	"github.com/dylanconnolly/drinkee/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) GenerateRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/drinks/:id", func(c *gin.Context) {
				// s.logger.Debug("info log in http server", map[string]string{
				// 	"route": c.FullPath(),
				// })
				s.logger.Debug(&logger.LogFields{
					Message: "debug log in http server",
					Fields: &logger.HttpFields{
						Method: c.Request.Method,
						Route:  c.Request.RequestURI,
					},
				})
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
