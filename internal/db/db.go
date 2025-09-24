package db

import (
	"fmt"
	"log"
	
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Failed to turn on foreign keys", err)
	}

	return db
}

func Migrate(db *sqlx.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	fmt.Println("migration applied successfully")
}
