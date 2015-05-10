package models

import (
	"fmt"
	"testing"
)

func TestNewUser(t *testing.T) {
	u := NewUser("test", "test123", "Test User")
	fmt.Println(u)

	users := LoadUsers("../data/test_users.json")
	users.AddUser(u)
	users.Save("../data/test_users.json")

	fmt.Println(*users)
}
