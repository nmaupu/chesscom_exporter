package main

import (
	"fmt"
	"github.com/nmaupu/chesscom_exporter/pkg/api/chesscom"
	"os"
	"time"
)

func test() {
	username := "nm0p"

	archives, err := chesscom.GetAllPlayerArchives(username)
	if err != nil {
		fmt.Printf("An error occurred getting player's archives for %s, err=%v", username, err)
	}

	file, err := os.OpenFile("/tmp/chesscom_games.pgn", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file for writing")
	}
	defer file.Close()

	for _, archive := range archives.Archives {
		games, err := chesscom.GetPlayerMonthlyArchivesByURL(archive.GetURL())
		if err != nil {
			fmt.Printf("An error occurred getting player's monthly archive for %s, err=%v", username, err)
		}

		// Saving PGN for all games
		for _, game := range games.Games {
			endTime := time.Unix(game.EndTime, 0)
			fmt.Printf("Processing game from %s: %s vs. %s", endTime, game.White.Username, game.Black.Username)
			_, err := file.Write([]byte(game.PGN + "\n"))
			if err != nil {
				fmt.Printf("Error writing to file, err=%v", err)
			}
		}
	}
}
