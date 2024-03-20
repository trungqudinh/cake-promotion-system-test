package domain

import (
	"cake/pkg/domain/errors"
	"cake/pkg/storage/mysql"
	"cake/pkg/time"
	"fmt"
)

type UserRepository interface {
	CreateUser(user *UserProfile) (uint, error)
	UpdateLastLogin(user *UserAccount) error
	CreateEvent(event *mysql.UserEvent) error
	FindUserByField(fieldName string, value string) (*UserAccount, error)
}

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

func (u UserStatus) String() string {
	switch u {
	case UserStatusDisable:
		return "disable"
	case UserStatusEnable:
		return "enable"
	default:
		return "unknown"
	}
}

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

func NewUserMySQLRepository(mysql *mysql.MySqlStorage) *UserMysqlRepository {
	fmt.Printf("NewUserMySQLRepository: %v\n", mysql)
	return &UserMysqlRepository{mysql: mysql}
}

type UserMysqlRepository struct {
	mysql *mysql.MySqlStorage
}

func (u *UserMysqlRepository) FindUserByField(fieldName string, value string) (*UserAccount, error) {
	return nil, nil
}

func (u *UserMysqlRepository) CreateUser(userProfile *UserProfile) (userId uint, err error) {
	user := mysql.User{
		FullName: userProfile.FullName,
		Username: userProfile.Username,
		Phone:    userProfile.Phone,
		Email:    userProfile.Email,
		Password: userProfile.Password,
	}
	user.Birthday, err = time.ParseLocal("2006-01-02", userProfile.Birthday)
	if err != nil {
		return
	}

	userIdentities := []mysql.UserIdentity{}

	if userProfile.Username != "" {
		userIdentities = append(userIdentities, mysql.UserIdentity{
			Type:  "username",
			Value: userProfile.Username,
		})
	}
	if userProfile.Phone != "" {
		userIdentities = append(userIdentities, mysql.UserIdentity{
			Type:  "phone",
			Value: userProfile.Phone,
		})
	}
	if userProfile.Email != "" {
		userIdentities = append(userIdentities, mysql.UserIdentity{
			Type:  "email",
			Value: userProfile.Email,
		})
	}

	userIdentity := &mysql.UserIdentity{}
	for _, i := range userIdentities {
		userIdentity, _ = u.mysql.FindUserIdentity(i.Value)
		if userIdentity != nil {
			err = errors.WrapResponse(
				nil,
				409,
				409,
				fmt.Sprintf("This %s is already registered", i.Type))
			return
		}
	}

	err = u.mysql.Create(&user)
	if err != nil {
		err = errors.Wrap500Response(err, "Create user failed")
		return
	}

	for _, i := range userIdentities {
		i.UserID = user.ID
		err = u.mysql.Create(&i)
		if err != nil {
			err = errors.Wrap500Response(err, "Create user identity failed %v", i)
			return
		}
	}
	return user.ID, nil
}

func (u *UserMysqlRepository) CreateEvent(event *mysql.UserEvent) error {
	return nil
}

func (u *UserMysqlRepository) UpdateLastLogin(user *UserAccount) error {
	return nil
}
