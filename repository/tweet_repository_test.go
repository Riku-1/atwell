package repository

import (
	"atwell/domain"
	"atwell/infrastructure"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	db, err := infrastructure.GetDevGormDB()
	if err != nil {
		t.Fatal(err)
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	r := NewMysqlTweetRepository(tx)
	testComment := "test_get"
	from := time.Now()
	tw := domain.Tweet{Comment: testComment}
	err = tx.Create(&tw).Error
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	twList, err := r.Get(from, from.AddDate(0, 0, 1))
	assert.NotNil(t, twList)
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
