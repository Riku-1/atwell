package domain

import (
	"time"

	"gorm.io/gorm"
)

type Tweet struct {
	gorm.Model
	UserID  uint
	Comment string
}

type TweetUsecase interface {
	Get(email string, from time.Time, to time.Time) ([]Tweet, error)
	Create(email string, comment string) (Tweet, error)
	Delete(email string, id uint) error
}

type TweetRepository interface {
	Get(user User, from time.Time, to time.Time) ([]Tweet, error)
	Create(user User, comment string) (Tweet, error)
	Delete(user User, id uint) error
}
