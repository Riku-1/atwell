package repository

import (
	"golang-api/domain"
	"gorm.io/gorm"
)

// mysqlArticleRepository is a article repository using mysql.
type mysqlArticleRepository struct {
	db *gorm.DB
}

// NewMysqlArticleRepository provides a mysqlArticleRepository struct.
func NewMysqlArticleRepository(db *gorm.DB) domain.ArticleRepository {
	return mysqlArticleRepository{db}
}

// GetAll returns all articles.
func (r mysqlArticleRepository) GetAll() (res []domain.Article, err error) {
	err = r.db.Find(&res).Error

	return
}
