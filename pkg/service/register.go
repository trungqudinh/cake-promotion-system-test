package service

import (
	"cake/pkg/api"
	"cake/pkg/convert"
	"cake/pkg/domain/errors"
	"net/mail"
	"regexp"
	"strings"
)

type RegisterService struct{}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Username string `json:"username" binding:"required_without_all=Phone Email"`
	Phone    string `json:"phone" binding:"required_without_all=Username Email"`
	Email    string `json:"email" binding:"required_without_all=Username Phone"`
	Password string `json:"password" binding:"required"`
	Birthday string `json:"birthday" binding:"required"`
}

type RegisterResponse struct {
	UserId   uint32 `json:"user_id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPhoneNumber(phoneNumber string) bool {
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")

	return re.Find([]byte(phoneNumber)) != nil
}

func IsValidUsername(username string) bool {
	usernameRegex := `^[a-zA-Z0-9_-]{3,16}$`
	re := regexp.MustCompile(usernameRegex)
	return re.Find([]byte(username)) != nil
}

func (r *RegisterService) ValidateUserIdentity(req *RegisterRequest) error {
	if req.Username != "" && !IsValidUsername(req.Username) {
		return errors.ErrInvalidUsername
	}
	if req.Phone != "" && !IsValidPhoneNumber(req.Phone) {
		return errors.ErrInvalidPhoneNumber
	}
	if req.Email != "" && !IsValidEmail(req.Email) {
		return errors.ErrInvalidEmail
	}
	return nil

}

func NewRegisterService() *RegisterService {
	return &RegisterService{}
}

func (r *RegisterService) Register(req *RegisterRequest) (response api.Response) {
	err := r.ValidateUserIdentity(req)
	if err != nil {
		response = api.Response{
			Status:  api.StatusBadRequest,
			Message: convert.ToPointer(err.Error()),
		}
		return
	}

	response = api.Response{
		Status:  api.StatusSuccess,
		Message: convert.ToPointer("Register Success"),
		Data: RegisterResponse{
			UserId:   1,
			FullName: req.FullName,
			Username: req.Username,
			Phone:    req.Phone,
			Email:    req.Email,
			Birthday: req.Birthday,
		},
	}
	return
}
