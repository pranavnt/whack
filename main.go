package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
	"github.com/muesli/termenv"
)

const (
	port = 23234

	rules = `Rules:

   • Clicking a target wins a point for your team

   • Clicking the other team's emoji loses a point and makes water

   • Clicking on water takes away a point`
)

var (
	b = NewBoard()

	fireScore = 0
	iceScore  = 0
)

func main() {
	rand.Seed(time.Now().UnixNano())
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf(":%d", port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			myCustomBubbleteaMiddleware(),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on port %d", port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

var programs = make([]*tea.Program, 0, 100)

var currTeam = true

func myCustomBubbleteaMiddleware() wish.Middleware {
	teaHandler := func(s ssh.Session) *tea.Program {
		m := &model{
			team: currTeam,
		}
		currTeam = !currTeam

		p := tea.NewProgram(m, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen(),
			tea.WithMouseCellMotion())

		m.thisProgram = p
		programs = append(programs, p)
		return p
	}
	return bm.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}

type model struct {
	team        bool // true means fire
	thisProgram *tea.Program
	x           int
	y           int

	comment string
}

func (m model) Init() tea.Cmd {
	return nil
}

// I tried not using this and it _looked_ like it still worked.
var lock = new(sync.Mutex)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// fmt.Println(fmt.Sprintf("%p", m.thisProgram), msg)
	switch msg := msg.(type) {
	//case tea.WindowSizeMsg:
	//	m.height = msg.Height
	//	m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			for i, p := range programs {
				if p == m.thisProgram {
					programs = append(programs[:i], programs[i+1:]...)
					break
				}
			}
			fmt.Println(m.thisProgram, "quitting")
			return m, tea.Quit
		}
	case tea.MouseMsg:
		if msg.Type != tea.MouseRelease { // trigger on release only - no dragging allowed
			return m, nil
		}
		m.x = (msg.X - 1) / 2 // divide by 2: each emoji is two cells wide
		m.y = msg.Y - 3       // subtract 2: the top two rows are not part of the board and the border isn't either

		m.comment = b.Click(m.x, m.y, m.team)

		lock.Lock()

		for _, p := range programs {
			if p == m.thisProgram {
				continue
			}
			p.Send(tea.Msg(true)) // trigger render
		}

		lock.Unlock()
	}

	return m, nil
}

var gameDoneMsg = ""

func (m model) View() string {
	if len(gameDoneMsg) > 0 {
		l := len([]rune(gameDoneMsg))
		return strings.Repeat("\n", height/2) + strings.Repeat(" ", width) + "Game over!\n" + strings.Repeat(" ", width-l/2) + gameDoneMsg + strings.Repeat(" ", width-l/2)
	}
	t := ""
	if m.team {
		t = "🔥"
	} else {
		t = "🧊"
	}
	return "You're in the " + t + " team! Click on targets to win " + t + "s for your team!\n" +
		"\n" +
		b.RenderBoard(t, fireScore, iceScore, m.comment) + "\n" + rules
}
