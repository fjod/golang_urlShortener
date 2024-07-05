package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Url struct {
	id  int64
	url string
}

//go:generate mockery --name Operations
type Operations interface {
	GetUrl(id int) (Url, error)
	SetUrlId() (int, error)
	SetUrl(string, int64) error
	CreateTable() error
}

type Service struct {
	Storage Operations
}

type SQLiteRepository struct {
	db *sql.DB
}

const fileName = "db/db.db"

func NewSQLiteRepository() *SQLiteRepository {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		fmt.Println("cant open db", err)
		panic(err)
	}
	return &SQLiteRepository{
		db: db,
	}
}

func NewService(storage Operations) *Service {
	return &Service{
		Storage: storage,
	}
}

func (r *SQLiteRepository) CreateTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS urls(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL        
    );`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) SetUrlId() (int64, error) {
	var rowId int64
	query := `INSERT INTO urls(url) VALUES (?) returning id`
	row := r.db.QueryRow(query, "")
	err := row.Scan(&rowId)
	return rowId, err
}

func (r *SQLiteRepository) SetUrl(url string, id int64) error {
	query := `update urls set url =? where id =?`
	_, err := r.db.Exec(query, url, id)
	return err
}

func (r *SQLiteRepository) GetUrl(id int) (Url, error) {
	var url Url
	query := `select url from urls where id =?`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&url.url)
	return url, err
}
