package config

import (
	"testing"
)

func TestGetPrdGormDB(t *testing.T) {
	_, err := GetPrdGormDB()

	if err != nil {
		t.Fatal(err)
	}
}
func TestGetDevGormDB(t *testing.T) {
	_, err := GetDevGormDB()

	if err != nil {
		t.Fatal(err)
	}
}
