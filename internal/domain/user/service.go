package user

type Service struct {
    repo Repository
}

func NewService(r Repository) *Service {
    return &Service{repo: r}
}

func (s *Service) RegisterUser(u *User) error {
    // business logic here
    return s.repo.Create(u)
}
