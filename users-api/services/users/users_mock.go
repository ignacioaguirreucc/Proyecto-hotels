package users

import (
	"users-api/domain/users"

	"github.com/stretchr/testify/mock"
)

// Mock struct para el Servicio
type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) GetAll() ([]users.User, error) {
	args := m.Called()
	if err := args.Error(1); err != nil {
		return nil, err
	}
	return args.Get(0).([]users.User), nil
}

func (m *Mock) GetByID(id int64) (users.User, error) {
	args := m.Called(id)
	if err := args.Error(1); err != nil {
		return users.User{}, err
	}
	return args.Get(0).(users.User), nil
}

func (m *Mock) Create(user users.User) (int64, error) {
	args := m.Called(user)
	if err := args.Error(1); err != nil {
		return 0, err
	}
	return args.Get(0).(int64), nil
}

func (m *Mock) Login(username string, password string) (users.LoginResponse, error) {
	args := m.Called(username, password)
	if err := args.Error(1); err != nil {
		return users.LoginResponse{}, err
	}
	return args.Get(0).(users.LoginResponse), nil
}
