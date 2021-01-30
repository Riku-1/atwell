package config

import (
	"testing"
)

func TestGetPrdDBConfig(t *testing.T) {
	_, err := GetTestDBConfig()

	if err != nil {
		t.Fatal(err)
	}
}
