package model

import (
	"database/sql"
	"fmt"
	"strings"
)

type PlayerScore struct {
	Id int
	Name string
	Score int
}

type row struct {
	id int
	name string
	score int
}


func (receiver PlayerScore) Delete(db *sql.DB) sql.Result  {
	stmt, err := db.Prepare("DELETE FROM player_score WHERE id = ?")
	checkErr(err)

	result, err := stmt.Exec(receiver.Id)
	checkErr(err)

	return result
}


func DeleteScores(db *sql.DB) {
	_, err := db.Exec("DELETE FROM player_score")
	checkErr(err)
}

func FindAllScores(db *sql.DB, limit int, offset int) []PlayerScore  {
	stmt, err := db.Prepare("SELECT id, name, score FROM player_score ORDER BY score DESC LIMIT ? OFFSET ?")
	checkErr(err)

	rowList, err := stmt.Query(limit, offset)
	checkErr(err)

	var result []PlayerScore
	for rowList.Next() {
		var r row
		rowList.Scan(&r.id, &r.name, &r.score)
		result = append(result, PlayerScore{Id: r.id, Name: r.name, Score: r.score})
	}

	return result
}

func FindScoreCount(db *sql.DB) int {
	rowList, err := db.Query("SELECT COUNT(*) FROM player_score")
	checkErr(err)

	var result int
	for rowList.Next() {
		rowList.Scan(&result)
	}

	return result
}

func FindScoreByName(name string, db *sql.DB) *PlayerScore  {
	stmt, err := db.Prepare("SELECT id, name, score FROM player_score WHERE name = ?")
	checkErr(err)
	rowList, err := stmt.Query(name)
	checkErr(err)

	var result PlayerScore
	for rowList.Next() {
		var r row
		rowList.Scan(&r.id, &r.name, &r.score)
		result = PlayerScore{Id: r.id, Name: r.name, Score: r.score}
	}

	return &result
}

func (receiver PlayerScore) Save(db *sql.DB) sql.Result {
	stmt, err := db.Prepare("INSERT INTO player_score (name, score) VALUES ($1, $2) ON CONFLICT(name) DO UPDATE SET score = $2 WHERE name = $1 AND score < $2")
	checkErr(err)

	result, err := stmt.Exec(receiver.Name, receiver.Score)
	checkErr(err)

	return result
}

func BulkSave(scoreList []PlayerScore, db *sql.DB) sql.Result {
	valuesString := make([]string, 0, len(scoreList))
	valuesArg := make([]interface{}, 0, len(scoreList) * 2)
	for _, score := range scoreList {
		valuesString = append(valuesString, "(?, ?)")
		valuesArg = append(valuesArg, score.Name)
		valuesArg = append(valuesArg, score.Score)
	}
	query := fmt.Sprintf("INSERT OR REPLACE INTO player_score (name, score) VALUES %s", strings.Join(valuesString, ","))
	stmt, err := db.Prepare(query)
	result, err := stmt.Exec(valuesArg...)
	checkErr(err)

	return result
}
