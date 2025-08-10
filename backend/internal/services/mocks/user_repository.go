package mocks

import "github.com/stretchr/testify/mock"

// MockUserRepository は UserRepositoryInterface のモック
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id string) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (interface{}, error) {
	args := m.Called(email)
	return args.Get(0), args.Error(1)
}

func (m *MockUserRepository) GetByGoogleID(googleID string) (interface{}, error) {
	args := m.Called(googleID)
	return args.Get(0), args.Error(1)
}

func (m *MockUserRepository) Create(user interface{}) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user interface{}) error {
	args := m.Called(user)
	return args.Error(0)
}
