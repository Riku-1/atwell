package domain

import (
	"time"

	"gorm.io/gorm"
)

type Tweet struct {
	gorm.Model
	Comment string
}

type TweetUsecase interface {
	Get(from time.Time, to time.Time) ([]Tweet, error)
	Create(comment string) (Tweet, error)
	Delete(id int) error
}

type TweetRepository interface {
	Get(from time.Time, to time.Time) ([]Tweet, error)
	Create(comment string) (Tweet, error)
	Delete(id int) error
}
