package domain

func GetGameStatus(board []int) string {
	winner := checkWinner(board)
	if winner == 1 {
		return "win"
	}
	if winner == 2 {
		return "lose"
	}
	if isBoardFull(board) {
		return "draw"
	}
	return "play"
}

func IsValidMove(oldBoard, newBoard []int) bool {
	newMovesCount := 0
	for i := 0; i < 9; i++ {
		if oldBoard[i] == newBoard[i] {
			continue
		}
		if oldBoard[i] == 0 && newBoard[i] == 1 {
			newMovesCount++
		} else {
			return false
		}
	}
	return newMovesCount == 1
}

func MakeServerMove(board []int) {
	bestScore := -1000
	bestMove := -1

	for i := 0; i < 9; i++ {
		if board[i] == 0 {
			board[i] = 2
			score := minimax(board, 0, false)
			board[i] = 0

			if score > bestScore {
				bestScore = score
				bestMove = i
			}
		}
	}

	if bestMove != -1 {
		board[bestMove] = 2
	}
}

func minimax(board []int, depth int, isMaximizing bool) int {
	winner := checkWinner(board)
	if winner == 2 {
		return 10 - depth
	}
	if winner == 1 {
		return depth - 10
	}
	if isBoardFull(board) {
		return 0
	}

	if isMaximizing {
		bestScore := -1000
		for i := 0; i < 9; i++ {
			if board[i] == 0 {
				board[i] = 2
				score := minimax(board, depth+1, false)
				board[i] = 0
				if score > bestScore {
					bestScore = score
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for i := 0; i < 9; i++ {
			if board[i] == 0 {
				board[i] = 1
				score := minimax(board, depth+1, true)
				board[i] = 0
				if score < bestScore {
					bestScore = score
				}
			}
		}
		return bestScore
	}
}

func checkWinner(b []int) int {
	lines := [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	for _, line := range lines {
		if b[line[0]] != 0 && b[line[0]] == b[line[1]] && b[line[0]] == b[line[2]] {
			return b[line[0]]
		}
	}
	return 0
}

func isBoardFull(b []int) bool {
	for _, cell := range b {
		if cell == 0 {
			return false
		}
	}
	return true
}
