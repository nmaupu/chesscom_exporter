package model

import (
	"strconv"
	"strings"
	"time"
)

type ChesscomArchive string

// ChesscomArchives represents an array of monthly games archive
type ChesscomArchives struct {
	Archives []ChesscomArchive `json:"archives"`
}

func (a ChesscomArchive) GetMonth() int {
	toks := strings.Split(string(a), "/")
	idx := len(toks) - 1
	if idx < 0 {
		return -1
	}
	month, err := strconv.Atoi(toks[idx])
	if err != nil {
		return -1
	}
	return month
}

func (a ChesscomArchive) GetMonthAsString() string {
	month := a.GetMonth()
	if month < 0 {
		return ""
	}
	t := time.Date(0, time.Month(month+1), 0, 0, 0, 0, 0, time.Local)
	return t.Month().String()
}

func (a ChesscomArchive) GetYear() int {
	toks := strings.Split(string(a), "/")
	idx := len(toks) - 2
	if idx < 0 {
		return -1
	}
	year, err := strconv.Atoi(toks[idx])
	if err != nil {
		return -1
	}
	return year
}

func (a ChesscomArchive) GetPlayerName() string {
	toks := strings.Split(string(a), "/")
	idx := len(toks) - 4
	if idx < 0 {
		return ""
	}
	return toks[idx]
}

func (a ChesscomArchive) GetURL() string {
	return string(a)
}
