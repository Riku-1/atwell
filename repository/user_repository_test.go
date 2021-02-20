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
	resUser, err := r.Create("test_create_user@email.com")

	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// assert response
	assert.Equal(t, "test_create_user@email.com", resUser.Email)
	tx.Commit()

	// confirm record created
	var user domain.User
	db.First(&user, resUser.ID)
	assert.Equal(t, "test_create_user@email.com", user.Email)

	// delete record
	db.Delete(&domain.User{}, user.ID)
}

func TestMysqlUserRepository_Create_WhenCrateDuplicateUser(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	r := NewMysqlUserRepository(tx)

	// create user
	_, err = r.Create("test_duplicate_user@email.com")
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	// create user by already exist email
	_, err = r.Create("test_duplicate_user@email.com")
	if err == nil {
		t.Fatal("duplicate error should occur")
	}
	assert.Error(t, err, infrastructure.DuplicateError{})
}
