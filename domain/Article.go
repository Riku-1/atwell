package domain

// Article is an interface of articles of blog.
type Article struct {
	Title       string `json:"Title"`
	Body        string `json:"body"`
	PublishDate string `json:"publish_date"`
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
