package service

import (
	"cake/pkg/domain"
	"cake/pkg/domain/errors"
	"cake/pkg/encrypt"
	"context"
)

type UserService struct {
	userRepo domain.UserRepository
}

func (u *UserService) Verify(ctx context.Context, identityField string, identityValue string, password string) (*domain.UserAccount, error) {
	user, err := u.userRepo.FindUserByField(identityField, identityValue)
	if err != nil {
		return nil, err
	}
	if user.Status != domain.UserStatusEnable {
		return nil, errors.ErrUserDisable
	}
	if user.Password != encrypt.SHA1String(password) {
		return nil, errors.ErrInvalidPassword
	}
	return user, nil
}

func (u *UserService) Login(ctx context.Context, identityField string, identityValue string, password string) (*domain.UserAccount, error) {
	user, err := u.Verify(ctx, identityField, identityValue, password)
	if err != nil {
		return nil, err
	}
	err = u.userRepo.UpdateLastLogin(user)
	if err != nil {
		return nil, err
	}
	return user, nil
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
