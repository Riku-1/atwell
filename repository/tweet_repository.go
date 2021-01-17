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
	err = r.db.Where("created_at BETWEEN ? AND ?", from, to).Find(&res).Error

	return
}
