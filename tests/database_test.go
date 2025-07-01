package tests

import (
	"testing"

	"GoServer/repository"
)

func TestConnection(t *testing.T) {
	var ret string

	err := repository.Database_connection.QueryRow("SELECT CURRENT_TIMESTAMP").Scan(&ret)
	if err != nil {
		t.Error(err)
	}
}
