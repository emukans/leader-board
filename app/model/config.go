package model

import (
	"strconv"
	"time"
)

type Config struct {
	Id int
	Name string
	Value string
	UpdatedAt time.Time
	CreatedAt time.Time
}


const (
	authToken = "auth_token"
	pageLimit = "page_limit"
)

func FindAuthToken() string {
	return FindConfigByName(authToken).Value
}

func FindPageLimit() (int, error) {
	return strconv.Atoi(FindConfigByName(pageLimit).Value)
}

func FindConfigByName(name string) *Config {
	stmt, err := DB.Prepare("SELECT id, name, value, updated_at, created_at FROM config WHERE name = ? LIMIT 1")
	checkErr(err)
	rowList, err := stmt.Query(name)
	checkErr(err)

	var result Config
	for rowList.Next() {
		rowList.Scan(&result.Id, &result.Name, &result.Value, &result.UpdatedAt, &result.CreatedAt)
	}

	return &result
}
