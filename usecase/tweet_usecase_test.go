package usecase

import (
	"atwell/domain"
	mocks "atwell/mocks/domain"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTweetUsecase_Get(t *testing.T) {
	// mock
	userRepository := new(mocks.UserRepository)
	email := "test_get@email.com"
	mockUser := domain.User{Email: email}
	userRepository.On("Get", email).Return(
		mockUser,
		nil,
	)

	tweetRepository := new(mocks.TweetRepository)
	mockedTwList := make([]domain.Tweet, 0)
	mockedTwList = append(
		mockedTwList,
		domain.Tweet{Comment: "test_get"},
	)
	from := time.Now()
	to := time.Now()
	tweetRepository.On("Get", mockUser, from, to).Return(mockedTwList, nil).Once()

	// call function
	u := NewTweetUsecase(
		tweetRepository,
		userRepository,
	)
	twList, err := u.Get(
		email,
		from,
		to,
	)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test_get", twList[0].Comment)
}

func TestTweetUsecase_Get_WhenUserRepositoryError(t *testing.T) {
	// mock
	userRepository := new(mocks.UserRepository)
	email := "test_get@email.com"
	mockUser := domain.User{Email: email}
	userRepository.On("Get", email).Return(
		mockUser,
		errors.New("some error"), // error occurred
	).Once()

	tweetRepository := new(mocks.TweetRepository)
	mockedTwList := make([]domain.Tweet, 0)
	mockedTwList = append(
		mockedTwList,
		domain.Tweet{Comment: "test_get"},
	)
	from := time.Now()
	to := time.Now()
	tweetRepository.On("Get", mockUser, from, to).Return(mockedTwList, nil).Once()

	// call function
	u := NewTweetUsecase(
		tweetRepository,
		userRepository,
	)
	_, err := u.Get(
		email,
		from,
		to,
	)
	if err == nil {
		t.Fatal("error should be returned when some error occurred")
	}
}

func TestTweetUsecase_Get_WhenTweetRepositoryError(t *testing.T) {
	// mock
	userRepository := new(mocks.UserRepository)
	email := "test_get@email.com"
	mockUser := domain.User{Email: email}
	userRepository.On("Get", email).Return(
		mockUser,
		nil,
	)

	tweetRepository := new(mocks.TweetRepository)
	mockedTwList := make([]domain.Tweet, 0)
	mockedTwList = append(
		mockedTwList,
		domain.Tweet{Comment: "test_get"},
	)
	from := time.Now()
	to := time.Now()
	tweetRepository.On("Get", mockUser, from, to).Return(
		mockedTwList,
		errors.New("some error"), // error occurred
	).Once()

	// call function
	u := NewTweetUsecase(
		tweetRepository,
		userRepository,
	)
	_, err := u.Get(
		email,
		from,
		to,
	)
	if err == nil {
		t.Fatal("error should be returned when some error occurred")
	}
}

//func TestCreate(t *testing.T) {
//	repo := new(mocks.TweetRepository)
//	comment := "test_create"
//	mockTweet := domain.Tweet{Comment: comment}
//
//	repo.On("Create", comment).Return(mockTweet, nil).Once()
//	u := NewTweetUsecase(repo)
//	tweet, err := u.Create(comment)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	assert.Equal(t, mockTweet.Comment, tweet.Comment)
//}
//
//func TestDelete(t *testing.T) {
//	repo := new(mocks.TweetRepository)
//	targetID := 111
//
//	repo.On("Delete", targetID).Return(nil).Once()
//	u := NewTweetUsecase(repo)
//	err := u.Delete(targetID)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
