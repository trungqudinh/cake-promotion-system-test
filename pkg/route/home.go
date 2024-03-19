package route

import (
	"cake/pkg/api"
	"cake/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Home(c *gin.Context) {
	HomeService := service.NewHomeService()
	api.JSON(c, http.StatusOK, HomeService.Home())
}
