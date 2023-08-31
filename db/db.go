package db

import (
	"Avito_Backend_Trainee/models"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	HOST = "database"
	//HOST = "localhost"
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

func NewNullTimeStamp(s string) sql.NullTime {
	if len(s) == 0 {
		return sql.NullTime{}
	}
	time, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  time,
		Valid: true,
	}
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
	ctx := context.Background()
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	query := `INSERT INTO segments (name) VALUES ($1);`
	_, err = tx.ExecContext(ctx, query, segment.Name)
	if err != nil {
		return err
	}

	if segment.Percent > 0 {
		var usersTotal float64
		query := `SELECT COUNT(*) as total FROM users`

		rows, err := tx.QueryContext(ctx, query)
		defer rows.Close()
		if err != nil {
			tx.Rollback()
			return err
		}
		for rows.Next() {
			err = rows.Scan(&usersTotal)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		usersCount := int(usersTotal * (float64(segment.Percent) / 100))
		query = `SELECT user_id FROM users ORDER BY random() LIMIT $1;`
		rows2, err := tx.QueryContext(ctx, query, usersCount)
		defer rows2.Close()
		if err != nil {
			tx.Rollback()
			return err
		}

		userList := []models.User{}
		for rows2.Next() {
			var user models.User
			err := rows2.Scan(&user.User_id)
			if err != nil {
				tx.Rollback()
				return err
			}
			userList = append(userList, models.User{User_id: user.User_id})
		}
		for _, user := range userList {
			updateUserSegment := models.UserSegmentsUpdate{
				User_id:               user.User_id,
				Segment_to_join_names: []models.Segment_to_join_name{{Name: segment.Name}},
			}
			for _, segToJoin := range updateUserSegment.Segment_to_join_names {
				query := `INSERT INTO userssegments (user_id, segment_name, deleted) VALUES ($1,$2,$3);`
				_, err := tx.ExecContext(ctx, query, updateUserSegment.User_id, segToJoin.Name, NewNullTimeStamp(segToJoin.Deleted))
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

	}
	err = tx.Commit()
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

	ctx := context.Background()
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, segToJDel := range userSegments.Segment_to_delete_names {
		query := `UPDATE userssegments SET deleted=$1 WHERE user_id=$2 AND segment_name=$3;`
		_, err := tx.ExecContext(ctx, query, time.Now(), userSegments.User_id, segToJDel.Name)
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	for _, segToJoin := range userSegments.Segment_to_join_names {
		query := `INSERT INTO userssegments (user_id, segment_name, deleted) VALUES ($1,$2,$3);`
		_, err := tx.ExecContext(ctx, query, userSegments.User_id, segToJoin.Name, NewNullTimeStamp(segToJoin.Deleted))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db Database) GetUserActiveSegments(user *models.User) (*models.UserAvailableSegmentList, error) {
	list := &models.UserAvailableSegmentList{}

	query := `SELECT user_id, segment_name FROM userssegments WHERE user_id=$1 and deleted IS NULL ORDER BY ID DESC;`
	rows, err := db.Conn.Query(query, user.User_id)
	defer rows.Close()
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

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (db Database) CreateReport(report *models.UserReport) error {
	ctx := context.Background()
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	query := `SELECT user_id, segment_name, joined, deleted FROM userssegments WHERE user_id=$1 AND joined>=$2 AND (deleted<=$3 or deleted IS NULL);`
	rows, err := tx.QueryContext(ctx, query, report.User_id, report.Start_date, report.End_date)
	defer rows.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	type History struct {
		userId      string
		segmentName string
		startTime   string
		endTime     sql.NullString
	}

	headers := []string{"userId", "segmentName", "startTime", "endTime"}
	records := []History{}
	for rows.Next() {
		var history History
		err := rows.Scan(&history.userId, &history.segmentName, &history.startTime, &history.endTime)
		if err != nil {
			tx.Rollback()
			return err
		}
		records = append(records, history)
	}

	if len(records) == 0 {
		tx.Rollback()
		return fmt.Errorf("no such user_id")
	}

	csvFileName := fmt.Sprintf("reports/%s_%s.csv", report.User_id, GenerateSecureToken(10))
	csvFile, err := os.Create(csvFileName)
	defer csvFile.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	if err := csvwriter.Write(headers); err != nil {
		tx.Rollback()
		return err
	}

	for _, record := range records {
		row := []string{record.userId, record.segmentName, record.startTime, record.endTime.String}
		if err := csvwriter.Write(row); err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	report.Report_name = csvFileName

	return nil
}

func (db Database) DownloadReportByName(reportName string) (string, string, error) {
	dst := fmt.Sprintf("reports/%s", reportName)
	file, err := os.ReadFile(dst)
	if err != nil {
		return "", "", err
	}
	contentType := http.DetectContentType(file[:512])
	filename := "reports/" + reportName
	return contentType, filename, nil
}
