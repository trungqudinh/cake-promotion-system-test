package route

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

func (server *Server) Listen(addr string) {
	fmt.Println("Listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
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
		v1.GET("/", s.Home)
	}
}

func (server *Server) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "End point is available. Version: " + os.Getenv("APP_VERSION"),
	})
}
