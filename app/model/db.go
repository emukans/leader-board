package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"path"
	"path/filepath"
	"runtime"
)

var db *sql.DB


func dbPath() string {
	_, baseDir, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(baseDir))

	return filepath.Dir(dir) + "/../db/leader_board.db"
}

func GetDBConnection() (*sql.DB, error) {
	if db != nil {
		err := db.Ping()
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	path := dbPath()

	var err error
	db, err = sql.Open("sqlite3", path)

	return db, err
}
