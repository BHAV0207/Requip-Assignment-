package user

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}


func (s *Service) Create(req CreateUserRequest) (*User, error) {

	user := &User{
		Name:              req.Name,
		Email:             req.Email,
		PrimaryMobile:     req.PrimaryMobile,
		SecondaryMobile:   req.SecondaryMobile,
		Aadhaar:           req.Aadhaar,
		PAN:               req.PAN,
		PlaceOfBirth:      req.PlaceOfBirth,
		CurrentAddress:    req.CurrentAddress,
		PermanentAddress:  req.PermanentAddress,
	}

	if err := s.repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}