package web

import (
	"tictactoe-vs-minmax/internal"

	"github.com/google/uuid"
)

func ToDomainBoard(req GameRequest) internal.Board {
	var b internal.Board
	for i := range 3 {
		for j := range 3 {
			b[i][j] = internal.Cell(req.Board[i][j])
		}
	}
	return b
}

func ToResponse(g internal.Game) GameResponse {
	var b [3][3]int
	for i := range 3 {
		for j := range 3 {
			b[i][j] = int(g.Board[i][j])
		}
	}
	return GameResponse{ID: g.ID.String(), Board: b}
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
