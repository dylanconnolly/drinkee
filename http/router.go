package http

import "github.com/gin-gonic/gin"

func (s *Server) GenerateRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/drinks", func(c *gin.Context) {
				s.handleGetDrinks(c)
			})
			v1.POST("/drinks", func(c *gin.Context) {
				s.handleCreateDrink(c)
			})
		}
	}
}
