package internal

type GameService struct {
	NextMove(game Game) (Game, error)
	Validate(prev, current Board) error
	IsGameOver(board Board) bool
}
