package usecase

import "golang-api/domain"

// tweetUsecase ...
type tweetUsecase struct {
	repository domain.TweetRepository
}

// NewTweetUsecase provides a tweetUsecase struct
func NewTweetUsecase(r domain.TweetRepository) domain.TweetUsecase {
	return tweetUsecase{r}
}

// Get ...
func (u tweetUsecase) Get() (res []domain.Tweet, err error) {
	res, err = u.repository.Get()

	if err != nil {
		//TODO
		return
	}

	return
}
