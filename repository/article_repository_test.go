package repository

import (
	"testing"
)

// TestGetAll tests GetAll function in normal system.
func TestGetAll(t *testing.T) {
	r := NewDummyArticleRepository()
	articles, _ := r.GetAll()

	if articles[0].Title != "aaa" {
		t.Fatal("Failed Test")
	}
}
