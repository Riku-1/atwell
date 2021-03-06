package usecase

import (
	"atwell/config"
	"atwell/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthenticationUsecase is a interface for yahoo japan authentication
type AuthenticationUsecase struct {
	getEmailUsecase domain.GetUserEmailUsecase
	userRepo        domain.UserRepository
}

// NewAuthenticationUsecase returns AuthenticationUsecase struct.
func NewAuthenticationUsecase(getEmailUsecase domain.GetUserEmailUsecase, userRepo domain.UserRepository) domain.AuthenticationUsecase {
	return &AuthenticationUsecase{
		getEmailUsecase: getEmailUsecase,
		userRepo:        userRepo,
	}
}

// SignUp gets user email address and creates account.
func (a *AuthenticationUsecase) SignUp(code string, nonce string) error {
	email, err := a.getEmailUsecase.GetEmail(code, nonce)
	if err != nil {
		return err
	}

	_, err = a.userRepo.Create(email)
	if err != nil {
		return err
	}

	return nil
}

// BeforeLogin returns token which contains nonce value.
func (a *AuthenticationUsecase) BeforeLogin(nonce string) (token string, err error) {
	return a.getEmailUsecase.BeforeLogin(nonce)
}

// Login gets user email address and returns auth Code.
func (a *AuthenticationUsecase) Login(code string, nonce string) (token string, err error) {
	email, err := a.getEmailUsecase.GetEmail(code, nonce)
	if err != nil {
		return "", err
	}

	user, err := a.userRepo.Get(email)
	if err != nil {
		return "", err
	}

	// return atwell Code for session
	atwellToken := jwt.New(jwt.SigningMethodHS256)
	claims := atwellToken.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // TODO

	c, _ := config.GetAppConfig() // TODO: constructor injection
	return atwellToken.SignedString([]byte(c.Secret))
}
