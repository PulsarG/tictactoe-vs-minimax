package internal

type GameService interface {
	NextMove(game Game) (Game, error)
	Validate(prev, current Board) error
	IsGameOver(board Board) bool
}
