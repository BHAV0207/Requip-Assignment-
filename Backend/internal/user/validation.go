package user

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
