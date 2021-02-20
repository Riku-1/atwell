package repository

import (
	"atwell/domain"
	"atwell/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMysqlUserRepository_Create(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

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
