package main

import (
	"math/rand"
)

type CellType rune

const (
	Empty  CellType = '🌲'
	Fire   CellType = '🔥'
	Ice    CellType = '🧊'
	Whack  CellType = '🎯'
	width  int      = 32
	height int      = 15
)

type Board struct {
	board  [][]CellType
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

	b := &Board{
		board: board,
	}

	b.Generate()

	return b
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

func (b *Board) Generate() {
	b.whackX = rand.Intn(width)
	b.whackY = rand.Intn(height)
	for i := range b.board {
		for j := range b.board[i] {
			b.board[i][j] = Empty
		}
	}
	b.board[b.whackY][b.whackX] = Whack
}

func (b *Board) Click(x, y int) {

	if x > width || y > height {
		return
	}

	if b.board[y][x] == Whack {
		b.board[y][x] = Empty
		b.Generate()
	}
}
