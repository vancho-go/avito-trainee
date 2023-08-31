package models

import (
	"fmt"
	"net/http"
)

type Segment struct {
	Name    string `json:"name"`
	Percent int    `json:"percent,omitempty"`
}

func (s *Segment) Bind(r *http.Request) error {
	if s.Name == "" {
		return fmt.Errorf("name is a required field")
	}
	if s.Percent < 0 || s.Percent > 100 {
		return fmt.Errorf("percent is out of range [1;100]")
	}
	return nil
}

func (s *Segment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
