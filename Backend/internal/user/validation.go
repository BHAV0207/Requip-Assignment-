package user

import "github.com/google/uuid"

func (s *Service) validateCreateRequest(req CreateUserRequest) error {

	exists, err := s.repository.ExistsByEmail(req.Email)
	if err != nil {
		return err
	}

	if exists {
		return ErrEmailAlreadyExists
	}

	exists, err = s.repository.ExistsByPAN(req.PAN)
	if err != nil {
		return err
	}

	if exists {
		return ErrPANAlreadyExists
	}

	exists, err = s.repository.ExistsByAadhaar(req.Aadhaar)
	if err != nil {
		return err
	}

	if exists {
		return ErrAadhaarAlreadyExists
	}

	exists, err = s.repository.ExistsByPrimaryMobile(req.PrimaryMobile)
	if err != nil {
		return err
	}

	if exists {
		return ErrMobileAlreadyExists
	}

	return nil
}

func (s *Service) validateUpdateRequest(id uuid.UUID, req UpdateUserRequest) error {

	if req.Email != nil {
		exists, err := s.repository.ExistsByEmail(*req.Email)
		if err != nil {
			return err
		}
		if exists {
			user, _ := s.repository.FindByID(id)
			if user != nil && user.Email != *req.Email {
				return ErrEmailAlreadyExists
			}
		}
	}

	if req.PrimaryMobile != nil {
		exists, err := s.repository.ExistsByPrimaryMobile(*req.PrimaryMobile)
		if err != nil {
			return err
		}
		if exists {
			user, _ := s.repository.FindByID(id)
			if user != nil && user.PrimaryMobile != *req.PrimaryMobile {
				return ErrMobileAlreadyExists
			}
		}
	}

	return nil
}
