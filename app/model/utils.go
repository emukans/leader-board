package model

import "database/sql"

func calculateScoreRange(below int, above int, total int) (int, int) {
	var lowerBound int
	var upperBound int

	limit := (total - 1) / 2

	if below > limit && above > limit {
		lowerBound = limit
		upperBound = total - limit - 1
	} else if below <= limit && above <= limit {
		lowerBound = below
		upperBound = above
	} else if below <= limit && above >= limit {
		lowerBound = below
		upperBound = total - 1 - below
	} else if below >= limit && above <= limit {
		upperBound = above
		lowerBound = total - 1 - above
	}

	return lowerBound, upperBound
}

func buildScoresFromSQLResult(rowList *sql.Rows) []PlayerScore {
	var result []PlayerScore
	for rowList.Next() {
		var score PlayerScore
		rowList.Scan(&score.Id, &score.Name, &score.Score, &score.UpdatedAt, &score.CreatedAt)
		result = append(result, score)
	}
	return result
}
