package usecase

import (
	"golang-api/domain"
	"time"
)

// tweetUsecase ...
type tweetUsecase struct {
	repository domain.TweetRepository
}

// NewTweetUsecase provides a tweetUsecase struct
func NewTweetUsecase(r domain.TweetRepository) domain.TweetUsecase {
	return tweetUsecase{r}
}

// Get ...
func (u tweetUsecase) Get(from time.Time, to time.Time) (res []domain.Tweet, err error) {
	res, err = u.repository.Get(from, to)

	return
}

// Create ...
func (u tweetUsecase) Create(comment string) (res domain.Tweet, err error) {
	res, err = u.repository.Create(comment)

	return
}
