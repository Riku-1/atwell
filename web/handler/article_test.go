package handler

import (
	"encoding/json"
	"golang-api/domain"
	mocks "golang-api/mocks/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGetArticles is a test for "/articles" requests in normal system.
func TestGetArticles(t *testing.T) {

	mockArticle := domain.Article{
		Title:       "aaa",
		Body:        "aaa",
		PublishDate: time.Date(2020, time.December, 20, 12, 00, 00, 0, time.UTC),
	}
	mockListArticle := make([]domain.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)

	mockUsecase := new(mocks.ArticleUsecase)
	mockUsecase.On("GetAll").Return(mockListArticle, nil)

	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	h := ArticleHandler{mockUsecase}
	h.getAllArticles(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Fatalf("Expected status code is 200, but actual is %v", status)
	}

	var resArticles []domain.Article
	err = json.Unmarshal(res.Body.Bytes(), &resArticles)
	if err != nil {
		t.Fatal(err)
	}

	if title := resArticles[0].Title; title != "aaa" {
		t.Fatalf("Expected title is aaa, but actual is %v", title)
	}
}
