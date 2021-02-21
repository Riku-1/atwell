package usecase

import (
	"atwell/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// authUsecase is struct for usecase about auth.
type authUsecase struct {
	repository domain.UserRepository
}

// NewAuthUsecase provides a authUsecase struct
func NewAuthUsecase(r domain.UserRepository) domain.UserUsecase {
	return authUsecase{repository: r}
}

// SignIn creates user account.
func (u authUsecase) SignIn(email string) (domain.User, error) {
	return u.repository.Create(email)
}

// Login creates and returns jwt token for session.
func (u authUsecase) Login(email string) (string, error) {
	user, err := u.repository.Get(email)
	if err != nil {
		return "", err
	}

	// return atwell token for session
	atwellToken := jwt.New(jwt.SigningMethodHS256)
	claims := atwellToken.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token, err := atwellToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return token, nil
}
