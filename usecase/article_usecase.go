package usecase

import (
	"golang-api/domain"
)

// articleUsecase ...
type articleUsecase struct {
	repository domain.ArticleRepository
}

// NewArticleUsecase provides a articleUsecase struct
func NewArticleUsecase(r domain.ArticleRepository) domain.ArticleUsecase {
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
