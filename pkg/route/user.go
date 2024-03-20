package route

import (
	"cake/pkg/api"
	"cake/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	ResgisterService *service.RegisterService
}

func (u *UserRoute) Register(c *gin.Context) {
	httpCode, request, response := api.BindRequest[service.RegisterRequest](c)
	if httpCode != http.StatusOK {
		api.JSON(c, httpCode, response)
		return
	}
	api.JSON(c, http.StatusOK, u.ResgisterService.Register(&request))
}

func (UserRoute) Login(c *gin.Context) {
	LoginService := service.NewLoginService()
	httpCode, request, response := api.BindRequest[service.LoginRequest](c)
	if httpCode != http.StatusOK {
		api.JSON(c, httpCode, response)
		return
	}
	api.JSON(c, http.StatusOK, LoginService.Login(&request))
}
