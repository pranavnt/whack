package main

type CellType rune

const (
	Empty CellType = '🌲'
	Fire CellType = '🔥'
	Ice CellType = '🧊'
	Whack CellType = '🎯'
	width int = 32
	height int =15
)

type Board struct {
	board [][]CellType
	whackX int
	whackY int
}

func NewBoard() *Board {
	board := make([][]CellType, height)
	for i := range board {
		board[i] = make([]CellType, width)
		for j := range board[i] {
			board[i][j] = Empty
		}
	}
	return &Board{
		board: board,
	}
}

func (b *Board) RenderBoard() string {
	var s string
	for _, row := range b.board {
		for _, cell := range row {
			s += string(cell)
		}
		s += "\n"
	}
	return s
}