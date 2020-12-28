package usecase

import (
	"golang-api/domain"
	mocks "golang-api/mocks/domain"
	"testing"
	"time"
)

// TestGetAll tests GetAll function in normal system.
func TestGetAll(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:       "aaa",
		Body:        "aaa",
		PublishDate: time.Date(2020, time.December, 20, 12, 00, 00, 0, time.UTC),
	}

	mockListArticle := make([]domain.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetAll").Return(mockListArticle, nil).Once()

		u := NewArticleUsecase(mockArticleRepo)
		articles, _ := u.GetAll()

		if articles[0].Title != "aaa" {
			t.Fatal("Failed Test")
		}
	})
}
