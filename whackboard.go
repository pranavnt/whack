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

	b.Seed()

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

func (b *Board) Seed() {
	b.whackX = rand.Intn(width)
	b.whackY = rand.Intn(height)

	for i := range b.board {
		for j := range b.board[i] {
			if i == b.whackY && j == b.whackX {
				b.board[i][j] = Whack
			} else {
				r := rand.Intn(20)
				// 5% chance of fire, 5% chance of ice
				if r == 0 {
					b.board[i][j] = Fire
				} else if r == 1 {
					b.board[i][j] = Ice
				} else {
					b.board[i][j] = Empty
				}
			}
		}
	}
}
