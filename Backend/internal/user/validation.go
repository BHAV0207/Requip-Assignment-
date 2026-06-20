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

	if req.PAN != nil {
		exists, err := s.repository.ExistsByPAN(*req.PAN)
		if err != nil {
			return err
		}
		if exists {
			user, _ := s.repository.FindByID(id)
			if user != nil && user.PAN != *req.PAN {
				return ErrPANAlreadyExists
			}
		}
	}

	if req.Aadhaar != nil {
		exists, err := s.repository.ExistsByAadhaar(*req.Aadhaar)
		if err != nil {
			return err
		}
		if exists {
			user, _ := s.repository.FindByID(id)
			if user != nil && user.Aadhaar != *req.Aadhaar {
				return ErrAadhaarAlreadyExists
			}
		}
	}

	return nil
}
