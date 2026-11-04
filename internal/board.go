package internal

type Cell int

const (
	Empty    Cell = 0
	Player   Cell = 1
	Computer Cell = 2
)

type Board [3][3]Cell
