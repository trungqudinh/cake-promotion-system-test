package service

import "cake/pkg/api"

type LoginService struct{}

func NewLoginService() *LoginService {
	return &LoginService{}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required_without_all=Phone Email"`
	Phone    string `json:"phone" binding:"required_without_all=Username Email"`
	Email    string `json:"email" binding:"required_without_all=Username Phone"`
	Password string `json:"password" binding:"required"`
}

func (r *LoginService) Login(req *LoginRequest) (response api.Response) {
	return
}
