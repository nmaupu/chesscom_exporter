package model

type ChesscomGames struct {
	Games []ChesscomGame `json:"games"`
}

type ChesscomGame struct {
	URL        string `json:"url"`
	PGN        string `json:"pgn"`
	EndTime    int64  `json:"end_time"`
	Rated      bool   `json:"rated"`
	Accuracies struct {
		White float64 `json:"white"`
		Black float64 `json:"black"`
	} `json:"accuracies"`
	TCN          string             `json:"tcn"`
	UUID         string             `json:"uuid"`
	InitialSetup string             `json:"initial_setup"`
	FEN          string             `json:"fen"`
	TimeClass    string             `json:"time_class"`
	Rules        string             `json:"rules"`
	White        ChesscomPlayerInfo `json:"white"`
	Black        ChesscomPlayerInfo `json:"black"`
}
