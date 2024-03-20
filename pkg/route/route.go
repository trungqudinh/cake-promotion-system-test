package route

import (
	"fmt"
	"log"

	"cake/pkg/api/middlewares"
	"cake/pkg/domain"
	"cake/pkg/service"
	"cake/pkg/storage/mysql"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	mysqlStorage *mysql.MySqlStorage
	Router       *gin.Engine
	routes       Routes
	services     Services
	repositories Repositories
}

type Routes struct {
	home HomeRoute
	user UserRoute
}

type Services struct {
	registerService *service.RegisterService
}

type Repositories struct {
	userRepository domain.UserRepository
}

func (server *Server) InitializeDatabases() {
	server.mysqlStorage = mysql.InitDatabase()
}

func (server *Server) Initialize() {
	server.InitializeDatabases()
	server.InitializeRepositories()
	server.InitializeServices()
	server.InitializeRoutes()
	server.InitializeApi()
}

func (s *Server) InitializeRepositories() {
	s.repositories = Repositories{
		userRepository: domain.NewUserMySQLRepository(s.mysqlStorage),
	}
}

func (server *Server) Listen(port string) {
	fmt.Println("Listening on: ", port)
	log.Fatal(server.Router.Run(fmt.Sprintf(":%s", port)))
}

func (s *Server) InitializeServices() {
	s.services = Services{
		registerService: service.NewRegisterService(
			service.WithUserRepository(s.repositories.userRepository),
		),
	}
}

func (s *Server) InitializeRoutes() {
	s.routes = Routes{
		home: HomeRoute{},
		user: UserRoute{
			ResgisterService: s.services.registerService,
		},
	}
}

func (s *Server) InitializeApi() {
	fmt.Printf("InitializeRoutes: %#v\n", s)
	s.Router = gin.Default()
	v1 := s.Router.Group("/rest/api/v1/")
	v1.Use(middlewares.PrometheusMiddleware)
	{
		v1.GET("/metrics", func(c *gin.Context) {
			handler := promhttp.Handler()
			handler.ServeHTTP(c.Writer, c.Request)
		}) // Home route

		v1.GET("/home", s.routes.home.Home)
		v1.POST("/register", s.routes.user.Register)
		v1.POST("/login", s.routes.user.Login)
	}
}
