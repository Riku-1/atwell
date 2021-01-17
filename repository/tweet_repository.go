package repository

import (
	"golang-api/domain"
	"time"

	"gorm.io/gorm"
)

// mysqlTweetRepository is a tweet repository using mysql.
type mysqlTweetRepository struct {
	db *gorm.DB
}

// NewMysqlTweetRepository provides a mysqlTweetRepository struct.
func NewMysqlTweetRepository(db *gorm.DB) domain.TweetRepository {
	return mysqlTweetRepository{db: db}
}

// Get returns tweets.
func (r mysqlTweetRepository) Get(from time.Time, to time.Time) (res []domain.Tweet, err error) {
	err = r.db.Where("created_at BETWEEN ? AND ?", from, to).Order("created_at desc").Find(&res).Error

	return
}

// Create creates new tweet.
func (r mysqlTweetRepository) Create(comment string) (tweet domain.Tweet, err error) {
	tweet = domain.Tweet{Comment: comment}
	err = r.db.Create(&tweet).Error

	return
}
