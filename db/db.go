package db

import (
	"Avito_Backend_Trainee/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	HOST = "database"
	PORT = 5432
)

// ErrNoMatch возвращается, если мы запрашиваем строку, которой несуществует
var ErrNoMatch = fmt.Errorf("no matching record")

type Database struct {
	Conn *sql.DB
}

func Initialize(username, password, database string) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Connected to database")
	return db, nil
}

func (db Database) AddUser(user *models.User) error {
	query := `INSERT INTO users (user_id) VALUES ($1);`
	_, err := db.Conn.Exec(query, user.User_id)
	if err != nil {
		return err
	}
	return nil
}

func (db Database) AddSegment(segment *models.Segment) error {
	query := `INSERT INTO segments (name) VALUES ($1);`
	_, err := db.Conn.Exec(query, segment.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db Database) DeleteSegment(segment *models.Segment) error {
	query := `DELETE FROM segments WHERE name = $1;`
	_, err := db.Conn.Exec(query, segment.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db Database) UpdateUserSegments(userSegments *models.UserSegmentsUpdate) error {
	for _, segToJDel := range userSegments.Segment_to_delete_names {
		query := `UPDATE userstosegments SET deleted=$1 WHERE user_id=$2 AND segment_name=$3;`
		_, err := db.Conn.Exec(query, time.Now(), userSegments.User_id, segToJDel)
		if err != nil {
			fmt.Errorf("%s", err)
		}
	}
	for _, segToJoin := range userSegments.Segment_to_join_names {
		if segToJoin.Deleted.IsZero() {
			query := `INSERT INTO userstosegments (user_id, segment_name) VALUES ($1,$2);`
			_, err := db.Conn.Exec(query, userSegments.User_id, segToJoin.Name)
			if err != nil {
				fmt.Errorf("%s", err)
			}
		} else {
			query := `INSERT INTO userstosegments (user_id, segment_name,deleted) VALUES ($1,$2,$3);`
			_, err := db.Conn.Exec(query, userSegments.User_id, segToJoin.Name, segToJoin.Deleted)
			if err != nil {
				fmt.Errorf("%s", err)
			}
		}
	}
	return nil
}

func (db Database) GetUserActiveSegments(user *models.User) (*models.UserAvailableSegmentList, error) {
	list := &models.UserAvailableSegmentList{}

	query := `SELECT user_id, segment_name FROM userstosegments WHERE user_id=$1 and deleted IS NULL ORDER BY ID DESC;`
	rows, err := db.Conn.Query(query, user.User_id)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var usertosegment models.UserAvailableSegment
		err := rows.Scan(&usertosegment.User_id, &usertosegment.Segment_name)
		if err != nil {
			return list, err
		}
		list.UserAvailableSegments = append(list.UserAvailableSegments, usertosegment)
	}
	return list, nil
}
