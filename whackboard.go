package main

import (
	"math/rand"
	"strings"
)

type CellType string

const (
	Tree     CellType = "🌳"
	TreeHot  CellType = "🌴"
	TreeCold CellType = "🌲"
	Fire     CellType = "🔥"
	Ice      CellType = "🧊"
	Whack    CellType = "🎯"
	Water    CellType = "️💧"
	width    int      = 32
	height   int      = 15
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
			r := rand.Intn(40)

			if r == 0 {
				board[i][j] = Fire
			} else if r == 1 {
				board[i][j] = Ice
			} else {
				board[i][j] = Tree
			}
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
	s += "╭" + strings.Repeat("─", width*2) + "╮" + "\n"

	for _, row := range b.board {
		s += "│"
		for _, cell := range row {
			s += string(cell)
		}
		s += "│\n"
	}

	s += "╰" + strings.Repeat("─", width*2) + "╯" + "\n"

	return s
}

func (b *Board) Generate() {
	b.whackX = rand.Intn(width)
	b.whackY = rand.Intn(height)
	b.board[b.whackY][b.whackX] = Whack
}

func (b *Board) Click(x, y int, team bool) {
	if x > width || y > height {
		return
	}

	if b.board[y][x] == Whack {
		if team {
			b.board[y][x] = Fire
			fireScore++
		} else {
			b.board[y][x] = Ice
			iceScore++
		}
		b.Generate()
	} else if b.board[y][x] == Fire {
		if !team {
			b.board[y][x] = Water
			iceScore--
		}
	} else if b.board[y][x] == Ice {
		if team {
			b.board[y][x] = Water
			fireScore--
		}
	} else if b.board[y][x] == Tree {
		if team {
			b.board[y][x] = TreeHot
		} else {
			b.board[y][x] = TreeCold
		}
	} else if b.board[y][x] == Water {
		if team {
			fireScore--
		} else {
			iceScore--
		}
		// b.board[y][x] = Tree
	}
}
