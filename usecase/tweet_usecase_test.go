package usecase

import (
	"atwell/domain"
	mocks "atwell/mocks/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	repo := new(mocks.TweetRepository)

	tweet := domain.Tweet{Comment: "test_get"}
	mockedTwList := make([]domain.Tweet, 0)
	mockedTwList = append(
		mockedTwList,
		domain.Tweet{Comment: "test_get"},
	)
	from := time.Now()
	to := time.Now()

	repo.On("Get", from, to).Return(mockedTwList, nil).Once()
	u := NewTweetUsecase(repo)
	twList, err := u.Get(
		from,
		to,
	)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, tweet.Comment, twList[0].Comment)
}

func TestCreate(t *testing.T) {
	repo := new(mocks.TweetRepository)
	comment := "test_create"
	mockTweet := domain.Tweet{Comment: comment}

	repo.On("Create", comment).Return(mockTweet, nil).Once()
	u := NewTweetUsecase(repo)
	tweet, err := u.Create(comment)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockTweet.Comment, tweet.Comment)
}

func TestDelete(t *testing.T) {
	repo := new(mocks.TweetRepository)
	targetID := 111

	repo.On("Delete", targetID).Return(nil).Once()
	u := NewTweetUsecase(repo)
	err := u.Delete(targetID)
	if err != nil {
		t.Fatal(err)
	}
}
