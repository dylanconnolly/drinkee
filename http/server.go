package http

import (
	"fmt"
	"net/http"
	"time"

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
	}

	s.Router.Use(cors.Default())
	s.Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s \"%s %s %s %d %s\" \"%s\" \"%s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	s.GenerateRoutes(s.Router)

	return s
}

func (s *Server) Serve() {
	s.Router.Run()
}
