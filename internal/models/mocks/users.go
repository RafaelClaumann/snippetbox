package mocks

import (
	"time"

	"snippetbox.claumann.net/internal/models"
)

var mockUser = &models.User{
	ID:             1,
	Name:           "string",
	Email:          "string",
	HashedPassword: []byte("teste"),
	Created:        time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) UpdatePassword(id int, current, new string) error {
	switch id {
	case 1:
		return models.ErrNoRecord
	case 2:
		return models.ErrInvalidCredentials
	default:
		return nil
	}
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
