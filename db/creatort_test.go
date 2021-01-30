package db

import (
	"atwell/config"
	"testing"
)

func TestCreateGormDB(t *testing.T) {
	dc, err := config.GetTestDBConfig()

	if err != nil {
		t.Fatal(err)
	}

	_, err = CreateGormDB(&dc)

	if err != nil {
		t.Error("DB Connection is failed. Please Check environment variables.")
		t.Fatal(err)
	}
}
