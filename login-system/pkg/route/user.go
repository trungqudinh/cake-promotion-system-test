package route

import (
	"cake/pkg/api"
	"cake/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	RegisterService *service.RegisterService
	LoginService    *service.LoginService
}

func (u *UserRoute) Register(c *gin.Context) {
	httpCode, request, response := api.BindRequest[service.RegisterRequest](c)
	if httpCode != http.StatusOK {
		api.JSON(c, httpCode, response)
		return
	}
	api.JSON(c, http.StatusOK, u.RegisterService.Register(&request))
}

func (u *UserRoute) Login(c *gin.Context) {
	httpCode, request, response := api.BindRequest[service.LoginRequest](c)
	if httpCode != http.StatusOK {
		api.JSON(c, httpCode, response)
		return
	}
	api.JSON(c, http.StatusOK, u.LoginService.Login(&request))
}
