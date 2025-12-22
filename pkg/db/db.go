package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
)	



func NewDBConnection(host, user, password, dbname string, port int) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}


type DB struct {
	conn *sql.DB
}

func (db *DB) Close() error {
	err := db.conn.Close()
	if err != nil {
		log.Println("Error closing database connection:", err)
	}
	return err
}