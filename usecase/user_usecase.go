package usecase

import (
	"atwell/domain"
)

// authUsecase is struct for usecase about auth.
type authUsecase struct {
	repository domain.UserRepository
}

// SignIn creates user account.
func (u authUsecase) SignIn(email string) (domain.User, error) {
	return u.repository.Create(email)
}

// NewAuthUsecase provides a authUsecase struct
func NewAuthUsecase(r domain.UserRepository) domain.UserUsecase {
	return authUsecase{repository: r}
}
