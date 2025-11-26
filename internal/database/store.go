package database

import (
	"database/sql"
	"devlink/internal/models"	
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func New(filepath string) (*Store, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	
	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func initSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			company TEXT,
			url TEXT UNIQUE,
			location TEXT,
			type TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS subscribers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE,
			preference TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SaveJob(job models.Job) (bool, error) {
	stmt, err := s.db.Prepare("INSERT OR IGNORE INTO jobs (title, company, url, location, type) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if job.Type == "" {
		job.Type = "new_grad"
	}

	res, err := stmt.Exec(job.Title, job.Company, job.URL, job.Location, job.Type)
	if err != nil {
		return false, err
	}

	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (s *Store) SaveSubscriber(sub models.Subscriber) error {
	_, err := s.db.Exec("INSERT OR REPLACE INTO subscribers (email, preference) VALUES (?, ?)", sub.Email, sub.Preference)
	return err
}

func (s *Store) GetAllSubscribers() ([]models.Subscriber, error) {
	rows, err := s.db.Query("SELECT email, preference FROM subscribers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscriber
	for rows.Next() {
		var sub models.Subscriber
		if err := rows.Scan(&sub.Email, &sub.Preference); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}