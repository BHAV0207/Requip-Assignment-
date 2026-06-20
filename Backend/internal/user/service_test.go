package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of IRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) ExistsByPAN(PAN string) (bool, error) {
	args := m.Called(PAN)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) ExistsByAadhaar(Aadhaar string) (bool, error) {
	args := m.Called(Aadhaar)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) FindByID(id uuid.UUID) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetAll(offset, limit int) ([]User, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockRepository) Update(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) Delete(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) Count() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

func TestService_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := CreateUserRequest{
		Name:          "John Doe",
		Email:         "john@example.com",
		PrimaryMobile: "1234567890",
		Aadhaar:       "123456789012",
		PAN:           "ABCDE1234F",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("ExistsByEmail", req.Email).Return(false, nil).Once()
		mockRepo.On("ExistsByPAN", req.PAN).Return(false, nil).Once()
		mockRepo.On("ExistsByAadhaar", req.Aadhaar).Return(false, nil).Once()
		mockRepo.On("Create", mock.Anything).Return(nil).Once()

		user, err := service.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Name, user.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo.On("ExistsByEmail", req.Email).Return(true, nil).Once()

		user, err := service.Create(req)

		assert.Error(t, err)
		assert.Equal(t, ErrEmailAlreadyExists, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetByID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedUser := &User{Name: "John Doe"}
		mockRepo.On("FindByID", userID).Return(expectedUser, nil).Once()

		user, err := service.GetByID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("FindByID", userID).Return(nil, ErrUserNotFound).Once()

		user, err := service.GetByID(userID)

		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}
