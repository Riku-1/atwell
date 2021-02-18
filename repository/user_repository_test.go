package repository

import (
	"atwell/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	tx := db.Begin()
	r := NewMysqlUserRepository(tx)
	resUser, err := r.Create("test_create_user")

	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assert response
	assert.Equal(t, "test_create_user", resUser.Email)
	tx.Commit()

	// confirm record created
	var user domain.User
	db.First(&user, resUser.ID)
	assert.Equal(t, "test_create_user", user.Email)

	// delete record
	db.Delete(&domain.User{}, user.ID)
}
