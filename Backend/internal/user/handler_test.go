package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of IService
type MockService struct {
	mock.Mock
}

func (m *MockService) Create(_ context.Context, req CreateUserRequest) (*User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockService) GetByID(_ context.Context, id uuid.UUID) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockService) GetAll(_ context.Context, page, limit int) ([]User, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]User), int64(args.Int(1)), args.Error(2)
}

func (m *MockService) Update(_ context.Context, id uuid.UUID, req UpdateUserRequest) (*User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockService) Delete(_ context.Context, id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := NewHandler(mockService)

	t.Run("success", func(t *testing.T) {
		req := CreateUserRequest{
			Name:             "John Doe",
			Email:            "john@example.com",
			PrimaryMobile:    "9876543210",
			Aadhaar:          "123456789012",
			PAN:              "ABCDE1234F",
			DateOfBirth:      "1990-01-01",
			PlaceOfBirth:     "Mumbai",
			CurrentAddress:   "123 Street",
			PermanentAddress: "123 Street",
		}
		expectedUser := &User{Name: req.Name, Email: req.Email}
		mockService.On("Create", req).Return(expectedUser, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body, _ := json.Marshal(req)
		c.Request = httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "User created successfully", res["message"])
		mockService.AssertExpectations(t)
	})
}

func TestHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	handler := NewHandler(mockService)
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedUser := &User{Name: "John Doe"}
		mockService.On("GetByID", userID).Return(expectedUser, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: userID.String()}}
		c.Request = httptest.NewRequest("GET", "/api/v1/users/"+userID.String(), nil)

		handler.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetByID", userID).Return(nil, ErrUserNotFound).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: userID.String()}}
		c.Request = httptest.NewRequest("GET", "/api/v1/users/"+userID.String(), nil)

		handler.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
