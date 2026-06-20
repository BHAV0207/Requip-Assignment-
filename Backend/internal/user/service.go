package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(req CreateUserRequest) (*User, error) {

	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	user := NewUser(req)

	if err := s.repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByID(id uuid.UUID) (*User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) GetAll(page, limit int) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	users, err := s.repository.GetAll(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repository.Count()
	if err != nil {
		return nil, 0, err
	}
	
	return users, total, nil
}

func (s *Service) Update(id uuid.UUID, req UpdateUserRequest) (*User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.validateUpdateRequest(id, req); err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.PrimaryMobile != nil {
		user.PrimaryMobile = *req.PrimaryMobile
	}
	if req.SecondaryMobile != nil {
		user.SecondaryMobile = *req.SecondaryMobile
	}
	if req.Aadhaar != nil {
		user.Aadhaar = *req.Aadhaar
	}
	if req.PAN != nil {
		user.PAN = *req.PAN
	}
	if req.DateOfBirth != nil {
		dob, _ := time.Parse("2006-01-02", *req.DateOfBirth)
		user.DateOfBirth = &dob
	}
	if req.PlaceOfBirth != nil {
		user.PlaceOfBirth = *req.PlaceOfBirth
	}
	if req.CurrentAddress != nil {
		user.CurrentAddress = *req.CurrentAddress
	}
	if req.PermanentAddress != nil {
		user.PermanentAddress = *req.PermanentAddress
	}

	if err := s.repository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(user)
}
