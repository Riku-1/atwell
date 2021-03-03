package usecase

import (
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
func (a *AuthenticationUsecase) SignUp(authInfo domain.AuthenticationInformation) error {
	email, err := a.getEmailUsecase.GetEmail(authInfo)
	if err != nil {
		return err
	}

	_, err = a.userRepo.Create(email)
	if err != nil {
		return err
	}

	return nil
}

// Login gets user email address and returns auth Token.
func (a *AuthenticationUsecase) Login(authInfo domain.AuthenticationInformation) (token string, err error) {
	email, err := a.getEmailUsecase.GetEmail(authInfo)
	if err != nil {
		return "", err
	}

	user, err := a.userRepo.Get(email)
	if err != nil {
		return "", err
	}

	// return atwell Token for session
	atwellToken := jwt.New(jwt.SigningMethodHS256)
	claims := atwellToken.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return atwellToken.SignedString([]byte("secret"))
}
