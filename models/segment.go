package models

import (
	"fmt"
	"net/http"
)

type Segment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *Segment) Bind(r *http.Request) error {
	if s.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	return nil
}

func (s *Segment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
