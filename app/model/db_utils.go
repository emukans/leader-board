package model

import (
	"path"
	"path/filepath"
	"runtime"
)

var (
	DBPath = RootPath() + "/../db/leader_board.db"
)


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
