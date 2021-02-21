package repository

import (
	"atwell/domain"
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
func (r mysqlTweetRepository) Get(user domain.User, from time.Time, to time.Time) ([]domain.Tweet, error) {
	var tweetList []domain.Tweet
	err := r.db.Where("user_id = ?", user.ID).Where("created_at BETWEEN ? AND ?", from, to).Order("created_at desc").Find(&tweetList).Error
	if err != nil {
		return nil, err
	}

	return tweetList, nil
}

// Create creates a new tweet.
func (r mysqlTweetRepository) Create(user domain.User, comment string) (tweet domain.Tweet, err error) {
	tweet = domain.Tweet{
		Comment: comment,
		UserID:  user.ID,
	}
	err = r.db.Create(&tweet).Error

	return
}

// Delete deletes a tweet specified by id.
func (r mysqlTweetRepository) Delete(id int) error {
	result := r.db.Delete(&domain.Tweet{}, id)
	return result.Error
}
