package model

type ChesscomPlayerInfo struct {
	Rating   int    `json:"rating"`
	Result   string `json:"result"`
	URL      string `json:"@id"`
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}
