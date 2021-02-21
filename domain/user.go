package domain

import "gorm.io/gorm"

// User is the model of a user.
type User struct {
	gorm.Model
	Email  string
	Tweets []Tweet
}

// UserRepository is a user repository interface.
type UserRepository interface {
	Get(email string) (User, error)
	Create(email string) (User, error)
}

// UserUsecase is a user usecase interface.
type UserUsecase interface {
	SignIn(email string) (User, error)
	Login(email string) (string, error)
}
