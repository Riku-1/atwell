package repository

import (
	"golang-api/domain"
	"time"
)

// dummyArticleRepository is a dummy repository and will be replaced by another one which use database.
type dummyArticleRepository struct {
}

// NewDummyArticleRepository provides a dummyArticleRepository struct.
func NewDummyArticleRepository() domain.ArticleRepository {
	return dummyArticleRepository{}
}

func (d dummyArticleRepository) GetAll() (res []domain.Article, err error) {
	res = []domain.Article{
		{
			Title:       "aaa",
			Body:        "aaa",
			PublishDate: time.Date(2020, time.December, 20, 12, 00, 00, 0, time.UTC),
		},
	}

	return
}
