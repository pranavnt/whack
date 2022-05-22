package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
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
	host = ""
	port = 23234
)

var b = NewBoard()

var programs = make(map[*tea.Program]int, 100)

func main() {
	rand.Seed(time.Now().UnixNano())
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
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
	log.Printf("Starting SSH server on %s:%d", host, port)
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

var currID int

// You can write your own custom bubbletea middleware that wraps tea.Program.
// Make sure you set the program input and output to ssh.Session.
func myCustomBubbleteaMiddleware() wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		return tea.NewProgram(m, opts...)
	}
	teaHandler := func(s ssh.Session) *tea.Program {
		//pty, _, active := s.Pty()
		//if !active {
		//	fmt.Println("no active terminal, skipping")
		//	_ = s.Exit(1)
		//	return nil
		//}
		currID++
		m := &model{
			team:    rand.Float64() > 0.5,
			modelID: currID,
			//term:   pty.Term,
			//width:  pty.Window.Width,
			//height: pty.Window.Height,
		}
		p := newProg(m, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen(), tea.WithMouseCellMotion())
		programs[p] = m.modelID
		return p
	}
	return bm.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	//term   string
	//width  int
	//height int
	team    bool // true means fire
	modelID int
	x       int
	y       int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//case tea.WindowSizeMsg:
	//	m.height = msg.Height
	//	m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.MouseMsg:
		if msg.Type != tea.MouseRelease { // trigger on release only - no dragging allowed
			return m, nil
		}
		m.x = msg.X / 2 // divide by 2: each emoji is two cells wide
		m.y = msg.Y - 2 // subtract 2: the top two rows are not part of the board
		//fmt.Println("mouse", m.x, m.y)

		b.Click(m.x, m.y, m.team)
		for p, id := range programs {
			fmt.Printf("%d %d\n", id, m.modelID)
			if id == m.modelID {
				continue
			}
			fmt.Println("rendering")
			p.Send(tea.Msg(true)) // trigger render
			fmt.Println("rendered")
		}
	}

	return m, nil
}

func (m model) View() string {
	t := ""
	if m.team {
		t = "🔥"
	} else {
		t = "🧊"
	}
	s := "You're in the %s team! Click on targets to win %ss for your team!\n"
	s += "Press 'q' to quit\n"
	s += b.RenderBoard()
	return fmt.Sprintf(s, t, t)
}