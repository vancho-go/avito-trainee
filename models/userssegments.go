package models

import (
	"fmt"
	"net/http"
	"time"
)

type Segment_to_join_name struct {
	Name    string `json:"name"`
	Deleted string `json:"deleted,omitempty"`
}

type Segment_to_delete_name struct {
	Name string `json:"name"`
}

type UserSegmentsUpdate struct {
	User_id                 string                   `json:"user_id"`
	Segment_to_join_names   []Segment_to_join_name   `json:"segment_to_join_names,omitempty"`
	Segment_to_delete_names []Segment_to_delete_name `json:"segment_to_delete_names,omitempty"`
}

func (us *UserSegmentsUpdate) Bind(r *http.Request) error {
	if us.User_id == "" {
		return fmt.Errorf("user_id is a required field")
	}

	for _, stjn := range us.Segment_to_join_names {
		if stjn.Deleted != "" {
			deleted, err := time.Parse("2006-01-02 15:04:05", stjn.Deleted)
			if err != nil {
				return fmt.Errorf("start_date parsing error")
			}
			if deleted.Before(time.Now()) {
				return fmt.Errorf("deleted value before now date")
			}
		}
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

type UserReport struct {
	User_id     string `json:"user_id"`
	Start_date  string `json:"start_date"`
	End_date    string `json:"end_date"`
	Report_name string `json:"report_name,omitempty"`
}

func (ur *UserReport) TimeComparator() error {
	startTime, err := time.Parse("2006-01-02 15:04:05", ur.Start_date)
	if err != nil {
		return err
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", ur.End_date)
	if err != nil {
		return err
	}
	if startTime.After(endTime) {
		return fmt.Errorf("start_time is after end_time")
	}
	return nil
}

func (ur *UserReport) Bind(r *http.Request) error {
	if ur.User_id == "" {
		return fmt.Errorf("user_id is a required field")
	}

	_, err := time.Parse("2006-01-02 15:04:05", ur.Start_date)
	if err != nil {
		return fmt.Errorf("start_date parsing error")
	}
	_, err = time.Parse("2006-01-02 15:04:05", ur.End_date)
	if err != nil {
		return fmt.Errorf("end_date parsing error")
	}
	if err := ur.TimeComparator(); err != nil {
		return err
	}
	return nil
}
func (*UserReport) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
