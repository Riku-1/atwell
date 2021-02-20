package repository

import (
	"atwell/domain"
	"atwell/infrastructure"
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

	if err != nil && !strings.Contains("Duplicate entry", err.Error()) {
		return domain.User{}, infrastructure.DuplicateError{}
	}

	return user, err
}
