package route

import (
	"fmt"
	"log"

	"cake/pkg/api/middlewares"
	"cake/pkg/storage/mysql"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (server *Server) ConnectDB() {
	mysql.InitDatabase()
}

func (server *Server) Initialize() {
	server.ConnectDB()
	server.InitializeRoutes()
}

func (server *Server) Listen(port string) {
	fmt.Println("Listening on: ", port)
	log.Fatal(server.Router.Run(":" + port))
}

func (s *Server) InitializeRoutes() {
	s.Router = gin.Default()
	v1 := s.Router.Group("/")
	v1.Use(middlewares.PrometheusMiddleware)
	{
		v1.GET("/metrics", func(c *gin.Context) {
			handler := promhttp.Handler()
			handler.ServeHTTP(c.Writer, c.Request)
		}) // Home route
		v1.GET("/home", s.Home)
	}
}
