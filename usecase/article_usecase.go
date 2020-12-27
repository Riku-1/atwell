package usecase

import (
	"golang-api/domain"
	"golang-api/repository"
)

// articleUsecase ...
type articleUsecase struct {
	repository domain.ArticleRepository
}

// NewArticleUsecase provides a articleUsecase struct
func NewArticleUsecase() domain.ArticleUsecase {
	r := repository.NewDummyArticleRepository()
	return articleUsecase{r}
}

func (a articleUsecase) GetAll() (res []domain.Article, err error) {
	res, err = a.repository.GetAll()

	if err != nil {
		// TODO
		return
	}

	return
}
