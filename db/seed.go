package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leader-board/app/model"
	"os"
)

func Seed() {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	seedPath := workingDir + "/db/seed/player_score.json"
	jsonFile, err := os.Open(seedPath)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var playerScoreList []model.PlayerScore

	json.Unmarshal(byteValue, &playerScoreList)
	for _, score := range playerScoreList {
		score.Save()
	}

	fmt.Printf("Successfully saved %d players\n", len(playerScoreList))
}
