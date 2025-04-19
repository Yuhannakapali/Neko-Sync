package userapp

import "nekosync/internal/domain/user"

type Service struct {
	repo user.Repository
}

func NewService(r user.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) RegisterUser(u *user.User) error {
	// business logic here
	return s.repo.Create(u)
}
