package repository

import (
	"golang-api/domain"
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
func (r mysqlTweetRepository) Get() (res []domain.Tweet, err error) {
	err = r.db.Find(&res).Error

	return
}
