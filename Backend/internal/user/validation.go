package user

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) validateCreateRequest(ctx context.Context, req CreateUserRequest) error {

	exists, err := s.repository.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if exists {
		return ErrEmailAlreadyExists
	}

	exists, err = s.repository.ExistsByPAN(ctx, req.PAN)
	if err != nil {
		return err
	}

	if exists {
		return ErrPANAlreadyExists
	}

	exists, err = s.repository.ExistsByAadhaar(ctx, req.Aadhaar)
	if err != nil {
		return err
	}

	if exists {
		return ErrAadhaarAlreadyExists
	}

	exists, err = s.repository.ExistsByPrimaryMobile(ctx, req.PrimaryMobile)
	if err != nil {
		return err
	}

	if exists {
		return ErrMobileAlreadyExists
	}

	return nil
}

func (s *Service) validateUpdateRequest(ctx context.Context, id uuid.UUID, req UpdateUserRequest) error {

	if req.Email != nil {
		exists, err := s.repository.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return err
		}
		if exists {
			user, _ := s.repository.FindByID(ctx, id)
			if user != nil && user.Email != *req.Email {
				return ErrEmailAlreadyExists
			}
		}
	}

	if req.PrimaryMobile != nil {
		exists, err := s.repository.ExistsByPrimaryMobile(ctx, *req.PrimaryMobile)
		if err != nil {
			return err
		}
		if exists {
			user, _ := s.repository.FindByID(ctx, id)
			if user != nil && user.PrimaryMobile != *req.PrimaryMobile {
				return ErrMobileAlreadyExists
			}
		}
	}

	return nil
}
