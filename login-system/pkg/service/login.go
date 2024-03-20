package service

import (
	"cake/pkg/api"
	"cake/pkg/convert"
	"cake/pkg/domain"
	"cake/pkg/domain/errors"
	"cake/pkg/encrypt"
)

type LoginService struct {
	userRepository domain.UserRepository
}

func NewLoginService(userRepository domain.UserRepository) *LoginService {
	return &LoginService{
		userRepository: userRepository,
	}
}

type LoginRequest struct {
	IdentityValue string `json:"identity_value" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserId    uint   `json:"user_id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	LastLogin string `json:"last_login"`
}

func (s *LoginService) Login(req *LoginRequest) (response api.Response) {
	user, err := s.userRepository.FindUserByIdentity(req.IdentityValue)
	defer func() {
		if err != nil {
			response = api.Response{
				Status:  api.StatusError,
				Message: convert.ToPointer(err.Error()),
			}
		}
	}()
	if err != nil {
		return
	}
	if user == nil {
		response = api.Response{
			Status:  api.StatusNotFound,
			Message: convert.ToPointer("User not found"),
		}
		return
	}

	if user.Password != encrypt.SHA1String(req.Password) {
		err = errors.ErrInvalidPassword
		return
	}

	response = api.Response{
		Status: api.StatusSuccess,
		Data: LoginResponse{
			UserId:    user.UserID,
			Username:  user.Username,
			Phone:     user.Phone,
			Email:     user.Email,
			LastLogin: user.LastLogin,
		},
	}

	s.userRepository.UpdateLastLogin(user.UserID)
	return
}
