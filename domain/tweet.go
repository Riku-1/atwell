package domain

import (
	"gorm.io/gorm"
)

type Tweet struct {
	gorm.Model
	Comment string
}

type TweetUsecase interface {
	Get() ([]Tweet, error)
}

type TweetRepository interface {
	Get() ([]Tweet, error)
}
