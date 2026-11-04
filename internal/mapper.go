package internal

import "github.com/google/uuid"

func ToDomain(m GameModel) (Game, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return Game{}, err
	}
	var board Board
	for i := range 3 {
		for j := range 3 {
			board[i][j] = Cell(m.Board[i][j])
		}
	}
	return Game{ID: id, Board: board}, nil
}
func ToModel(g Game) GameModel {
	var board [3][3]int
	for i := range 3 {
		for j := range 3 {
			board[i][j] = int(g.Board[i][j])
		}
	}
	return GameModel{ID: g.ID.String(), Board: board}
}
