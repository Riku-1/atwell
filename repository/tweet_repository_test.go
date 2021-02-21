package repository

import (
	"atwell/domain"
	"atwell/infrastructure"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMysqlTweetRepository_Get(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
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
	twList, err := r.Get(user, time.Now().AddDate(0, 0, -1), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "tweet1", twList[0].Comment)
	assert.Equal(t, "tweet2", twList[1].Comment)
}

func TestCreate(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	r := NewMysqlTweetRepository(tx)

	testComment := "test_creeate"
	tweet, err := r.Create(testComment)

	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	assert.Equal(t, tweet.Comment, testComment)
}

func TestDelete(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	r := NewMysqlTweetRepository(tx)
	tweet := domain.Tweet{Comment: "test_delete"}
	err = tx.Create(&tweet).Error
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	err = r.Delete(int(tweet.ID))
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	var tweet2 domain.Tweet
	tx.Find(&tweet2, tweet.ID)

	assert.Equal(t, 0, int(tweet2.ID))
}
