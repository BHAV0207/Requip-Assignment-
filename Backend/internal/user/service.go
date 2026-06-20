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

	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	user := NewUser(req)

	if err := s.repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
