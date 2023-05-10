package http

import (
	"net/http"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/dylanconnolly/drinkee/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	server       *http.Server
	Router       *gin.Engine
	DrinkService drinkee.DrinkService
	logger       logger.Logger
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		Router: gin.New(),
		logger: logger.New(),
	}

	s.Router.Use(cors.Default())

	s.GenerateRoutes(s.Router)

	return s
}

func (s *Server) Serve() {
	s.Router.Run()
}
