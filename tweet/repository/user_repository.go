package repository

import (
	"atwell/domain"
	"atwell/infrastructure/db"
	"strings"

	"gorm.io/gorm"
)

// mysqlUserRepository is a user repository using mysql.
type mysqlUserRepository struct {
	db *gorm.DB
}

// NewMysqlUserRepository provides a mysqlUserRepository struct.
func NewMysqlUserRepository(db *gorm.DB) domain.UserRepository {
	return mysqlUserRepository{db: db}
}

// Create creates a new user.
// TODO: Verify email address is valid form
func (r mysqlUserRepository) Create(email string) (domain.User, error) {
	user := domain.User{Email: email}
	err := r.db.Create(&user).Error

	if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
		return domain.User{}, db.DuplicateError{}
	}

	return user, err
}

// Get finds a user by email.
func (r mysqlUserRepository) Get(email string) (domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil && strings.Contains(err.Error(), "record not found") {
		return domain.User{}, db.NotFoundError{}
	}

	return user, nil
}
