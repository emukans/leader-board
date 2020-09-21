package model

import (
	"database/sql"
	"path"
	"path/filepath"
	"runtime"
)

var (
	DBPath = RootPath() + "/../db/leader_board.db"
	DB = CreateConnection()
)

func CreateConnection() *sql.DB {
	db, err := sql.Open("sqlite3", DBPath)
	checkErr(err)

	return db
}


func RootPath() string {
	_, baseDir, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(baseDir))

	return filepath.Dir(dir)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
