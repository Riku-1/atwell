package usecase

import (
	"atwell/domain"
	mocks "atwell/mocks/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_SignIn(t *testing.T) {
	repo := new(mocks.UserRepository)
	email := "test_sign_in_user"
	user := domain.User{Email: email}
	repo.On("Create", email).Return(user, nil).Once()

	u := NewAuthUsecase(repo)
	user, err := u.SignIn(email)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, email, user.Email)
}
