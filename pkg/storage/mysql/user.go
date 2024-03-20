package mysql

import (
	"errors"
	"slices"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName string `gorm:"size:255;"`
	Username string `gorm:"size:255;"`
	Phone    string `gorm:"size:100;"`
	Email    string `gorm:"size:100;"`
	Password string `gorm:"size:100;not null;"`
	Birthday time.Time
}

type UserEvent struct {
	Id        uint `gorm:"primaryKey"`
	UserID    uint
	User      User
	Status    int `gorm:"default:0"`
	LastLogin time.Time
}

type UserIdentity struct {
	UserID uint
	User   User
	Type   string
	Value  string `gorm:"not null;"`
	Active bool   `gorm:"default:true"`
}

func (m *MySqlStorage) EnableDebug() *MySqlStorage {
	m.db = m.db.Debug()
	return m
}
func (m *MySqlStorage) FirstOrCreate(value interface{}) error {
	return m.db.FirstOrCreate(value).Error
}

func (m *MySqlStorage) Create(value interface{}) error {
	return m.db.Create(value).Error
}

func (m *MySqlStorage) FindOne(query *gorm.DB, out interface{}) (bool, error) {
	result := query.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *MySqlStorage) FindUserById(id uint) (*User, error) {
	var (
		err error
		u   User
	)
	query := m.db.Model(User{}).Where("id = ?", id)
	found, err := m.FindOne(query, &u)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &u, nil
}

func (m *MySqlStorage) FindUserIdentity(value string) (*UserIdentity, error) {
	var (
		err    error
		record UserIdentity
	)
	query := m.db.Model(UserIdentity{}).Where("active = 1 and value = ?", value)
	found, err := m.FindOne(query, &record)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &record, nil
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

func (m *MySqlStorage) GetLatestUserEvent(userId uint) (*UserEvent, error) {
	var (
		err error
		ue  UserEvent
	)
	err = m.db.Model(UserEvent{}).Where("user_id = ?", userId).Last(&ue).Error
	if err != nil {
		return nil, err
	}
	return &ue, nil
}
