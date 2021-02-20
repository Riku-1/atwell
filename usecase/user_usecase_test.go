package usecase

import (
	"atwell/domain"
	mocks "atwell/mocks/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_SignIn(t *testing.T) {
	repo := new(mocks.UserRepository)
	email := "test_auth_usecase_sign_in@email.com"
	user := domain.User{Email: email}
	repo.On("Create", email).Return(user, nil).Once()
	u := NewAuthUsecase(repo)

	user, err := u.SignIn(email)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, email, user.Email)
}

func TestAuthUsecase_Login(t *testing.T) {
	repo := new(mocks.UserRepository)
	email := "test_auth_usecase_login@email.com"
	user := domain.User{Email: email}
	repo.On("Get", email).Return(user, nil).Once()
	u := NewAuthUsecase(repo)

	token, err := u.Login(email)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, "", token)
}

func TestAuthUsecase_Login_WhenErrorOccurs(t *testing.T) {
	repo := new(mocks.UserRepository)
	email := "test_auth_usecase_login@email.com"
	user := domain.User{Email: email}
	returnedError := errors.New("some error") // some error occurs
	repo.On("Get", email).Return(user, returnedError).Once()
	u := NewAuthUsecase(repo)

	token, err := u.Login(email)
	assert.Error(t, err)
	assert.Equal(t, "", token)
}
