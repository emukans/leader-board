package model

import (
	"database/sql"
	"time"
)

type PlayerScore struct {
	Id        int       `json:"-"`
	Name      string    `json:"name"`
	Score     int       `json:"score"`
	Rank      int       `json:"rank"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

func (receiver PlayerScore) Delete() error {
	db, err := GetDBConnection()
	if err != nil {
		return err
	}
	var stmt *sql.Stmt
	stmt, err = db.Prepare("DELETE FROM player_score WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(receiver.Id)

	return err
}

func DeleteScores() error {
	db, err := GetDBConnection()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM player_score")

	return err
}

func FindAllScores(limit int, offset int, periodFrom time.Time) ([]PlayerScore, error) {
	var rowList *sql.Rows
	var err error
	var stmt *sql.Stmt

	db, err := GetDBConnection()
	if err != nil {
		return nil, err
	}

	if !periodFrom.IsZero() {
		stmt, err = db.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score WHERE updated_at >= ? GROUP BY name ORDER BY score DESC LIMIT ? OFFSET ?")
		if err != nil {
			return nil, err
		}

		rowList, err = stmt.Query(periodFrom, limit, offset)
		if err != nil {
			return nil, err
		}
	} else {
		stmt, err = db.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score GROUP BY name ORDER BY score DESC LIMIT ? OFFSET ?")
		if err != nil {
			return nil, err
		}

		rowList, err = stmt.Query(limit, offset)
		if err != nil {
			return nil, err
		}
	}

	result := buildScoresFromSQLResult(rowList)

	return result, nil
}

func FindScoresAroundName(name string, pageLimit int) ([]PlayerScore, error) {
	lookupScore, err := FindScoreByName(name)

	if *lookupScore == (PlayerScore{}) || err != nil {
		return nil, err
	}
	lookupScore.Rank, err = lookupScore.FindRank()
	if err != nil {
		return nil, err
	}

	var rowList *sql.Rows
	var stmt *sql.Stmt

	db, err := GetDBConnection()
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score WHERE score <= ? AND name != ? GROUP BY name ORDER BY score DESC LIMIT ?")
	if err != nil {
		return nil, err
	}

	rowList, err = stmt.Query(lookupScore.Score, name, pageLimit)
	if err != nil {
		return nil, err
	}

	scoresBelowLookup := buildScoresFromSQLResult(rowList)
	for i, _ := range scoresBelowLookup {
		scoresBelowLookup[i].Rank = lookupScore.Rank + i + 1
	}

	stmt, err = db.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score WHERE score > ? AND name != ? GROUP BY name ORDER BY score DESC LIMIT ?")
	if err != nil {
		return nil, err
	}

	rowList, err = stmt.Query(lookupScore.Score, name, pageLimit)
	if err != nil {
		return nil, err
	}

	scoresAboveLookup := buildScoresFromSQLResult(rowList)
	for i, _ := range scoresAboveLookup {
		scoresAboveLookup[i].Rank = lookupScore.Rank - len(scoresAboveLookup) + i
	}
	var result []PlayerScore
	scoresBelowCount := len(scoresBelowLookup)
	scoresAboveCount := len(scoresAboveLookup)

	lowerBound, upperBound := calculateScoreRange(scoresBelowCount, scoresAboveCount, pageLimit)
	result = append(scoresAboveLookup[len(scoresAboveLookup)-upperBound:], *lookupScore)
	result = append(result, scoresBelowLookup[:lowerBound]...)

	return result, nil
}

func FindScoreCount() (int, error) {
	db, err := GetDBConnection()
	if err != nil {
		return 0, err
	}

	rowList, err := db.Query("SELECT COUNT(DISTINCT name) FROM player_score")
	if err != nil {
		return 0, err
	}

	var result int
	for rowList.Next() {
		rowList.Scan(&result)
	}

	return result, nil
}

func FindScoreByName(name string) (*PlayerScore, error) {
	db, err := GetDBConnection()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare("SELECT id, name, score, updated_at, created_at FROM player_score WHERE name = ? ORDER BY score DESC LIMIT 1")
	if err != nil {
		return nil, err
	}

	rowList, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}

	var result PlayerScore
	for rowList.Next() {
		rowList.Scan(&result.Id, &result.Name, &result.Score, &result.UpdatedAt, &result.CreatedAt)
	}

	return &result, nil
}

func (receiver PlayerScore) Save() error {
	db, err := GetDBConnection()
	if err != nil {
		return err
	}

	if receiver.UpdatedAt.IsZero() {
		stmt, err := db.Prepare("INSERT INTO player_score (name, score) VALUES ($1, $2)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(receiver.Name, receiver.Score)

		return err
	} else {
		stmt, err := db.Prepare("INSERT INTO player_score (name, score, updated_at) VALUES ($1, $2, $3)")
		if err != nil {
			return err
		}

		_, err = stmt.Exec(receiver.Name, receiver.Score, receiver.UpdatedAt)

		return err
	}
}

func (receiver PlayerScore) FindRank() (int, error) {
	db, err := GetDBConnection()
	if err != nil {
		return 0, err
	}

	stmt, err := db.Prepare("SELECT COUNT(DISTINCT name) FROM player_score WHERE score > ? AND name != ?")
	if err != nil {
		return 0, err
	}

	rowList, err := stmt.Query(receiver.Score, receiver.Name)

	var result int
	for rowList.Next() {
		rowList.Scan(&result)
	}

	return result + 1, nil
}
