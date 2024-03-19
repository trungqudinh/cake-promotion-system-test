package mysql

import (
	"errors"
	"slices"
	"time"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	FullName  string    `gorm:"size:255;" json:"full_name"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Phone     string    `gorm:"size:100;not null;unique" json:"phone"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Birthday  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserEvent struct {
	UserID    User
	Status    int       `gorm:"default:0" json:"status"`
	Username  string    `gorm:"size:255;not null;" json:"username"`
	LastLogin time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_login"`
}

func (m *MySqlStorage) FirstOrCreate(value interface{}) error {
	return m.db.Debug().FirstOrCreate(value).Error
}

// Supported Fields: `{"id", "username", "phone", "email"}`
func (m *MySqlStorage) FindUserByField(fieldName string, value string) (*User, error) {
	var (
		err error
		u   User
	)
	supportedFields := []string{"id", "username", "phone", "email"}
	if !slices.Contains(supportedFields, fieldName) {
		return nil, errors.New("unsupported field: " + fieldName)
	}

	err = m.db.Model(User{}).Where(fieldName+" = ?", value).Take(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (m *MySqlStorage) GetLatestUserEvent(username string) (*UserEvent, error) {
	var (
		err error
		ue  UserEvent
	)
	err = m.db.Model(UserEvent{}).Where("username = ?", username).Last(&ue).Error
	if err != nil {
		return nil, err
	}
	return &ue, nil
}
