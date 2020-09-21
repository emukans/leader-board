package model

import (
	"database/sql"
	"time"
)

type PlayerScore struct {
	Id int `json:"-"`
	Name string `json:"name"`
	Score int `json:"score"`
	Rank int `json:"rank"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}


func (receiver PlayerScore) Delete() sql.Result  {
	stmt, err := DB.Prepare("DELETE FROM player_score WHERE id = ?")
	checkErr(err)

	result, err := stmt.Exec(receiver.Id)
	checkErr(err)

	return result
}


func DeleteScores() {
	_, err := DB.Exec("DELETE FROM player_score")
	checkErr(err)
}

func FindAllScores(limit int, offset int, periodFrom time.Time) []PlayerScore  {
	var rowList *sql.Rows
	var err error
	var stmt *sql.Stmt

	if !periodFrom.IsZero() {
		stmt, err = DB.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score WHERE updated_at >= ? ORDER BY score DESC LIMIT ? OFFSET ?")
		checkErr(err)

		rowList, err = stmt.Query(periodFrom, limit, offset)
		checkErr(err)
	} else {
		stmt, err = DB.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score ORDER BY score DESC LIMIT ? OFFSET ?")
		checkErr(err)

		rowList, err = stmt.Query(limit, offset)
		checkErr(err)
	}

	var result []PlayerScore
	for rowList.Next() {
		var score PlayerScore
		rowList.Scan(&score.Id, &score.Name, &score.Score, &score.UpdatedAt, &score.CreatedAt)
		result = append(result, score)
	}

	return result
}

func FindScoreCount() int {
	rowList, err := DB.Query("SELECT COUNT(*) FROM player_score")
	checkErr(err)

	var result int
	for rowList.Next() {
		rowList.Scan(&result)
	}

	return result
}

func FindScoreByName(name string) *PlayerScore  {
	stmt, err := DB.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score WHERE name = ?")
	checkErr(err)
	rowList, err := stmt.Query(name)
	checkErr(err)

	var result PlayerScore
	for rowList.Next() {
		rowList.Scan(&result.Id, &result.Name, &result.Score, &result.UpdatedAt, &result.CreatedAt)
	}

	return &result
}

func (receiver PlayerScore) Save() sql.Result {
	if receiver.UpdatedAt.IsZero() {
		stmt, err := DB.Prepare("INSERT INTO player_score (name, score) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET score = $2 WHERE name = $1 AND score < $2")
		checkErr(err)

		result, err := stmt.Exec(receiver.Name, receiver.Score)
		checkErr(err)

		return result
	} else {
		stmt, err := DB.Prepare("INSERT INTO player_score (name, score, updated_at) VALUES ($1, $2, $3)")
		checkErr(err)

		result, err := stmt.Exec(receiver.Name, receiver.Score, receiver.UpdatedAt)
		checkErr(err)

		return result
	}
}
