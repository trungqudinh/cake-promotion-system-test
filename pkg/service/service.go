package service

import (
	"cake/pkg/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

type UserServiceOption func(*UserService)

func NewUserService(opts ...UserServiceOption) *UserService {
	us := &UserService{}
	for _, opt := range opts {
		opt(us)
	}
	return us
}

func WithUserRepo(repo domain.UserRepository) UserServiceOption {
	return func(us *UserService) {
		us.userRepo = repo
	}
}
