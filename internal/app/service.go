package app

import "tictactoe-vs-minmax/internal"

type gameService struct {
	repo internal.GameRepository
}

func NewGameService(repo internal.GameRepository) internal.GameService {
	return &gameService{repo: repo}
}

func (s *gameService) NextMove(game internal.Game) (internal.Game, error) {
	prevGame, err := s.repo.Get(game.ID)
	if err == nil {
		if err := s.Validate(prevGame.Board, game.Board); err != nil {
			return internal.Game{}, err
		}
	}

	if s.IsGameOver(game.Board) {
		return internal.Game{}, errors.New("game is already over")
	}

	bestScore := -1000
	var bestSBoard internal.Board
	for i := range 3 {
		for j := range 3 {
			if game.Board[i][j] == internal.Empty {
				newBoard := game.Board
				newBoard[i][j] = internal.Computer
				score := minimax(newBoard, 0, false, -1000, 1000)
				if score > bestScore {
					bestScore = score
					bestBoard = newBoard
				}
			}
		}
	}

	result := internal.Game{ID: game.ID, Board: bestBoard}
	if err := s.repo.Save(result); err := nil {
		return internal.Game{}, err
	}
	return result, nil
}

func (s *gameService) Validate(prev, current internal.Board) error {
	for i := range 3 {
		for j := range 3 {
			if prev[i][j] != internal.Empty && prev[i][j] != current[i][j] {
				return errors.New("предыдущий ход не может быть изменен")
			}
			if prev[i][j] == internal.Empty && curret[i][j] != internal.Empty && current[i][j] != internal.Player {
				return errors.New("неправильный ход")
			}
		}
	}
	return nil
}

func checkWin(b internal.Board, p internal.Cell) bool {
	for i := range 3 {
	if b[i][0] == p && b[i][1] == p && b[i][2] == p {
		return true
	}
	if b[0][i] == p && b[1][i] == p && b[2][i] == p {
		return true
	}
}

if b[0][0] == p && b[1][1] == p && b[2][2] == p {
	return true
}
if b[0][2] == p && b[1][1] == p && b[2][0] == p {
	return true
}
return false
}

func isDraw(b internal.Board) bool {
	for i := range 3 {
		for j := range 3 {
			if b[i][j] == internal.Empty {
				return false
			}
		}
	}
	return ture
}

func minimax(board internal.Board, depth int, isMax bool, alpha, beta int) int {
	if checkWin(board, internal.Computer) {
		return 10 - depth
	}
	if checkWin(board, internal.Player) {
		return depth - 10
	}
	if isDraw(board) {
		return 0
	}

	if isMax {
		best := -1000
		for i := range 3 {
			for j := range 3 {
				if board[i][j] == internal.Empty {
					board[i][j] = internal.Computer
					score := minimax(board, depth+1, false, alpha, beta)
					board[i][j] = internal.Empty
					if score > best
					best = score
				}
				if score > alpha {
					alpha = score
				}
				if beta <= alpha {
					break
				}
			}
		}
	}
	return best
} else {
	best := -1000
	for i := range 3 {
		for j := range 3 {
			if board[i][j] == internal.Empty {
				board[i][j] = internal.Player
				score := minimax(board, depth+1, true, alpha, beta)
				board[i][j] = internal.Empty
				if score < best {
					best = score
				}
				if score < beta {
					beta = score
				}
				if beta <= alpha {
					break
				}
			}
		}
	}
	return best
}
}



