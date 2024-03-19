package service

import (
	"cake/config"
	"cake/pkg/api"
)

type HomeService struct{}

func NewHomeService() *HomeService {
	return &HomeService{}
}

func (h HomeService) Home() api.Response {
	return api.Response{
		Status: api.StatusSuccess,
		Data:   config.GetAppConfig(),
	}
}
