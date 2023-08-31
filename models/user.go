package models

import (
	"fmt"
	"net/http"
)

type User struct {
	User_id string `json:"user_id"`
}

func (u *User) Bind(r *http.Request) error {
	if u.User_id == "" {
		return fmt.Errorf("user_id is a required field")
	}
	return nil
}

func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
