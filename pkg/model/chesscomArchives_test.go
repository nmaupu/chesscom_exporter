package model

import "testing"

func TestChesscomArchive_GetMonth(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "test 1",
			url:  "https://api.chess.com/pub/player/erik/games/2007/07",
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ChesscomArchive(tt.url)
			if got := a.GetMonth(); got != tt.want {
				t.Errorf("GetMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChesscomArchive_GetPlayerName(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "test 1",
			url:  "https://api.chess.com/pub/player/erik/games/2007/07",
			want: "erik",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ChesscomArchive(tt.url)
			if got := a.GetPlayerName(); got != tt.want {
				t.Errorf("GetPlayerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChesscomArchive_GetYear(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "test 1",
			url:  "https://api.chess.com/pub/player/erik/games/2007/07",
			want: 2007,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ChesscomArchive(tt.url)
			if got := a.GetYear(); got != tt.want {
				t.Errorf("GetYear() = %v, want %v", got, tt.want)
			}
		})
	}
}
