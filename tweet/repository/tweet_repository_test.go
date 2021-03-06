package repository

import (
	"atwell/config"
	"atwell/domain"
	db2 "atwell/infrastructure/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMysqlTweetRepository_Get(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	// create user and tweet
	user := domain.User{
		Email: "mysql_tweet_repository_test_get@email.com",
		Tweets: []domain.Tweet{
			{Comment: "tweet1"},
			{Comment: "tweet2"},
		},
	}
	tx.Create(&user)

	r := NewMysqlTweetRepository(tx)
	twList, err := r.Get(user, time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, 7))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "tweet1", twList[0].Comment)
	assert.Equal(t, "tweet2", twList[1].Comment)
}

func TestCreate(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	// create user
	user := domain.User{
		Email: "mysql_tweet_repository_test_create@email.com",
	}
	tx.Create(&user)

	r := NewMysqlTweetRepository(tx)
	tweet, err := r.Create(user, "tweet_repository_create_test")

	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	assert.Equal(t, tweet.Comment, "tweet_repository_create_test")
}

func TestMysqlTweetRepository_Delete(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	// create user and tweet
	user := domain.User{
		Email:  "mysql_tweet_repository_test_delete@email.com",
		Tweets: []domain.Tweet{},
	}
	tx.Create(&user)

	tweet := domain.Tweet{
		Comment: "tweet_repository_test_delete",
		UserID:  user.ID,
	}
	tx.Create(&tweet)

	// call function
	r := NewMysqlTweetRepository(tx)
	err = r.Delete(user, tweet.ID)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	var deletedTweet domain.Tweet
	tx.Find(&deletedTweet)
	assert.NotNil(t, deletedTweet.DeletedAt)
}

func TestMysqlTweetRepository_Delete_ByNotOwner(t *testing.T) {
	db, err := config.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	// create user1 and tweet
	user1 := domain.User{
		Email:  "mysql_tweet_repository_test_delete@email.com",
		Tweets: []domain.Tweet{},
	}
	tx.Create(&user1)

	tweet := domain.Tweet{
		Comment: "tweet_repository_test_delete",
		UserID:  user1.ID,
	}
	tx.Create(&tweet)

	user2 := domain.User{
		Email:  "another_user@email.com",
		Tweets: []domain.Tweet{},
	}

	// call function by not owner
	r := NewMysqlTweetRepository(tx)
	err = r.Delete(user2, tweet.ID)

	assert.IsType(t, db2.NoAuthorizationError{}, err)
}
