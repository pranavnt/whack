package main

import (
	"math/rand"
	"strconv"
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
	width    int      = 15
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
			r := rand.Float64()

			if r < 0.02 {
				board[i][j] = Fire
			} else if r < 0.04 {
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

//func intTo3DigStr(i int) string {
//	if i <= -10 {
//		return strconv.Itoa(i)
//	} else if i < 0 {
//		return strconv.Itoa(i) + "─"
//	} else if i < 10 {
//		return strconv.Itoa(i) + "──"
//	} else if i < 100 {
//		return strconv.Itoa(i) + "─"
//	} else {
//		return strconv.Itoa(i)
//	}
//}

func (b *Board) RenderBoard(t string, fireScore, iceScore int, comment string) string {
	var s string

	//scoreStr := "🔥 " + strconv.Itoa(fireScore) + " 🧊 " + strconv.Itoa(iceScore)

	border := "─"
	scoreStr := "🔥 " + strconv.Itoa(fireScore) + " 🧊 " + strconv.Itoa(iceScore)
	//fmt.Println(len(scoreStr), len([]rune(scoreStr)))
	//fireStr := "🔥 " + strconv.Itoa(fireScore)
	//iceStr := "🧊 " + strconv.Itoa(iceScore)
	l := len([]rune(scoreStr))
	s += t + strings.Repeat(border, width-l/2-l%2-2) + scoreStr + strings.Repeat(border, width-l/2-2) + t + "\n"
	//s += "╭" + strings.Repeat("─", 4) + t + t + t + strings.Repeat("─", (width-13)*2) + "🔥 " + intTo3DigStr(fireScore) + "──" + "🧊 " + intTo3DigStr(iceScore) + "──" + "╮" + "\n"

	for _, row := range b.board {
		s += "│"
		for _, cell := range row {
			s += string(cell)
		}
		s += "│\n"
	}

	l = len([]rune(comment))
	s += t + strings.Repeat(border, width-l/2-l%2-1) + comment + strings.Repeat(border, width-l/2-1) + t + "\n"
	//s += t + strings.Repeat(border, 2*width-2) + t + "\n"

	return s
}

func (b *Board) Generate() {
	b.whackX = rand.Intn(width)
	b.whackY = rand.Intn(height)
	if b.board[b.whackY][b.whackX] == Tree || b.board[b.whackY][b.whackX] == TreeHot || b.board[b.whackY][b.whackX] == TreeCold {
		b.board[b.whackY][b.whackX] = Whack
		return
	}
	b.Generate()
}

func (b *Board) Click(x, y int, team bool) string {
	if x >= width || y >= height || x < 0 || y < 0 {
		return "Out of bounds!"
	}

	comment := ""

	if b.board[y][x] == Whack {
		if team {
			b.board[y][x] = Fire
			fireScore++
		} else {
			b.board[y][x] = Ice
			iceScore++
		}
		comment = "Nice!"
		b.Generate()
	} else if b.board[y][x] == Fire {
		if !team {
			b.board[y][x] = Water
			iceScore--
		}
		comment = "Ouch! You made water!"
	} else if b.board[y][x] == Ice {
		if team {
			b.board[y][x] = Water
			fireScore--
		}
		comment = "Ouch! You made water!"
	} else if b.board[y][x] == Tree {
		if team {
			b.board[y][x] = TreeHot
			comment = "Palm tree!"
		} else {
			b.board[y][x] = TreeCold
			comment = "Pine tree!"
		}
	} else if b.board[y][x] == Water {
		if team {
			fireScore--
			comment = "Ouch! Water puts out fire!"
		} else {
			iceScore--
			comment = "Ouch! Water melts ice!"
		}
		// b.board[y][x] = Tree
	}

out:
	for r, row := range b.board {
		for c, cell := range row {
			if cell == TreeHot || cell == TreeCold || cell == Tree {
				break out
			}

			if r == len(b.board)-1 && c == len(row)-1 {
				s := ""
				if fireScore > iceScore {
					s += "🔥 wins!"
				} else if iceScore > fireScore {
					s += "🧊 wins!"
				} else {
					s += "It's a tie!"
				}
				gameDoneMsg = s
				return ""
			}
		}
	}

	return comment
}
