package usecase

import (
	"atwell/domain"
	"time"
)

// tweetUsecase ...
type tweetUsecase struct {
	tweetRepository domain.TweetRepository
	userRepository  domain.UserRepository
}

// NewTweetUsecase provides a tweetUsecase struct
func NewTweetUsecase(r domain.TweetRepository, ur domain.UserRepository) domain.TweetUsecase {
	return tweetUsecase{r, ur}
}

// Get returns tweet list.
func (u tweetUsecase) Get(email string, from time.Time, to time.Time) ([]domain.Tweet, error) {
	user, err := u.userRepository.Get(email)
	if err != nil {
		return nil, err
	}

	return u.tweetRepository.Get(user, from, to)
}

// Create ...
func (u tweetUsecase) Create(email string, comment string) (domain.Tweet, error) {
	user, err := u.userRepository.Get(email)
	if err != nil {
		return domain.Tweet{}, err
	}

	return u.tweetRepository.Create(user, comment)
}

// Delete ...
func (u tweetUsecase) Delete(id int) error {
	return u.tweetRepository.Delete(id)
}
