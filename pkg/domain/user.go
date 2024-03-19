package domain

import "context"

type UserAccount struct {
	UserID    uint32     `json:"user_id"`
	Username  string     `json:"username"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	LastLogin string     `json:"last_login"`
	Status    UserStatus `json:"status"`
}

type UserStatus int

const (
	UserStatusDisable UserStatus = -1
	UserStatusEnable  UserStatus = 0
)

type UserProfile struct {
	UserAccount
	FullName  string `json:"full_name"`
	Birthday  string `json:"birthday"`
	CreatedAt string `json:"created_at"`
}

type UserRepository interface {
	FindUserByField(ctx context.Context, fieldName string, value string) (*UserAccount, error)
	CreateUser(ctx context.Context, user *UserProfile) error
	UpdateLastLogin(ctx context.Context, user *UserAccount) error
}
