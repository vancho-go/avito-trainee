package models

import (
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ID      int    `json:"id"`
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

type Segment_to_join_name struct {
	Name    string    `json:"name"`
	Deleted time.Time `json:"deleted,omitempty"`
}

type UserSegmentsUpdate struct {
	ID                      int                    `json:"id"`
	User_id                 string                 `json:"user_id"`
	Segment_to_join_names   []Segment_to_join_name `json:"segment_to_join_names,omitempty"`
	Segment_to_delete_names []string               `json:"segment_to_delete_names,omitempty"`
}

func (us *UserSegmentsUpdate) Bind(r *http.Request) error {
	if us.User_id == "" {
		return fmt.Errorf("user_id is a required field")
	}
	return nil
}

func (*UserSegmentsUpdate) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UserAvailableSegment struct {
	User_id      string `json:"user_id"`
	Segment_name string `json:"segment_name"`
}

func (uts *UserAvailableSegment) Bind(r *http.Request) error {
	if uts.User_id == "" {
		return fmt.Errorf("user_id is a required field")
	}
	if uts.Segment_name == "" {
		return fmt.Errorf("segment_id is a required field")
	}
	return nil
}

func (*UserAvailableSegment) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UserAvailableSegmentList struct {
	UserAvailableSegments []UserAvailableSegment `json:"user_available_segments"`
}

func (*UserAvailableSegmentList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
