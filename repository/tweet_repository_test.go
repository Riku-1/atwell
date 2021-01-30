package repository

import (
	"atwell/config"
	adb "atwell/db"
	"atwell/domain"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
)

var db *gorm.DB

func setup() {
	dc, err := config.GetTestDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err = adb.CreateGormDB(&dc)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func TestGet(t *testing.T) {
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	r := NewMysqlTweetRepository(tx)
	testComment := "test_get"
	from := time.Now()
	tw := domain.Tweet{Comment: testComment}
	err := tx.Create(&tw).Error
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	twList, err := r.Get(from, from.AddDate(0, 0, 1))
	assert.NotNil(t, twList)
}

func TestCreate(t *testing.T) {
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
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	r := NewMysqlTweetRepository(tx)
	tweet := domain.Tweet{Comment: "test_delete"}
	err := tx.Create(&tweet).Error
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
