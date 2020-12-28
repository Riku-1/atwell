package domain

import "time"

// Article is an interface of articles of blog.
type Article struct {
	Title       string    `json:"Title"`
	Body        string    `json:"body"`
	PublishDate time.Time `json:"publish_date"`
}

// ArticleUsecase ...
type ArticleUsecase interface {
	// GetAll returns all articles.
	GetAll() ([]Article, error)
}

// ArticleRepository ...
type ArticleRepository interface {
	// GetAll returns all articles.
	GetAll() ([]Article, error)
}
