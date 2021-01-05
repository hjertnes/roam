// Package selectinput a option-select widget based on bubbletea
package selectinput

import (
	"github.com/hjertnes/roam/models"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hjertnes/roam/errs"
	"github.com/rotisserie/eris"
)

type model struct {
	choices []models.Choice
	label   string
	cursor  int
	choice  chan models.Choice
}

const (
	zero = 0
	one  = 1
)

// Run displays the select and returns the result.
func Run(label string, choices []models.Choice) (*models.Choice, error) {
	result := make(chan models.Choice, 1)

	switch len(choices) {
	case zero:
		return nil, eris.Wrap(errs.ErrNotFound, "no choices found")
	case one:
		result <- choices[0]
	default:
		p := tea.NewProgram(initialModel(label, choices, result))
		if err := p.Start(); err != nil {
			return nil, eris.Wrap(err, "failed to get user selection")
		}
	}

	r := <-result

	return &r, nil
}



func initialModel(label string, choices []models.Choice, choice chan models.Choice) model {
	return model{
		label:   label,
		choices: choices,
		cursor:  0,
		choice:  choice,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			os.Exit(0)
		case "enter":
			m.choice <- m.choices[m.cursor]

			return m, tea.Quit
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(m.label)
	s.WriteString("\n")

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("[x] ")
		} else {
			s.WriteString("[ ] ")
		}

		s.WriteString(m.choices[i].Title)
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
