package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dylanconnolly/drinkee/drinkee"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	server       *http.Server
	Router       *gin.Engine
	DrinkService drinkee.DrinkService
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		Router: gin.New(),
	}

	s.Router.Use(cors.Default())
	s.Router.Use(gin.Recovery())
	s.SetLogOutputDest()
	s.ApplyLogFormat()

	s.GenerateRoutes(s.Router)

	return s
}

func (s *Server) Serve() {
	s.Router.Run()
}

func (s *Server) SetLogOutputDest() {
	file, _ := os.Create("router.log")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}

func (s *Server) ApplyLogFormat() {
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
}
