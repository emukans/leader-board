package model

import (
	"strconv"
	"time"
)

type Config struct {
	Id        int
	Name      string
	Value     string
	UpdatedAt time.Time
	CreatedAt time.Time
}

const (
	authToken = "auth_token"
	pageLimit = "page_limit"
)

func FindAuthToken() (string, error) {
	config, err := FindConfigByName(authToken)
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

func FindPageLimit() (int, error) {
	config, err := FindConfigByName(pageLimit)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(config.Value)
}

func FindConfigByName(name string) (*Config, error) {
	db, err := GetDBConnection()
	if err != nil {
		return nil, err
	}
	stmt, err := db.Prepare("SELECT id, name, value, updated_at, created_at FROM config WHERE name = ? LIMIT 1")

	if err != nil {
		println(err)
		return nil, err
	}

	rowList, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	var result Config

	for rowList.Next() {
		rowList.Scan(&result.Id, &result.Name, &result.Value, &result.UpdatedAt, &result.CreatedAt)
	}

	return &result, nil
}
