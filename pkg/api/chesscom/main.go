package chesscom

import (
	"encoding/json"
	"fmt"
	"github.com/nmaupu/chesscom_exporter/pkg/model"
	"io/ioutil"
	"net/http"
)

const (
	chesscomAPI = "https://api.chess.com/pub"
)

func GetAllPlayerArchives(username string) (*model.ChesscomArchives, error) {
	url := fmt.Sprintf("%s/player/%s/games/archives", chesscomAPI, username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	archives := model.ChesscomArchives{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &archives); err != nil {
		return nil, err
	}

	return &archives, nil
}

func GetPlayerMonthlyArchives(username string, year int, month int) (*model.ChesscomGames, error) {
	url := fmt.Sprintf("%s/player/%s/games/%d/%d", chesscomAPI, username, year, month)
	return GetPlayerMonthlyArchivesByURL(url)
}

func GetPlayerMonthlyArchivesByURL(url string) (*model.ChesscomGames, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	games := model.ChesscomGames{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &games); err != nil {
		return nil, err
	}

	return &games, nil
}
