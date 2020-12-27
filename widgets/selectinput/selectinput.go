package selectinput

import (
	"strings"


	tea "github.com/charmbracelet/bubbletea"
	"github.com/hjertnes/roam/errs"
	"github.com/rotisserie/eris"
)

type model struct {
	choices []Choice
	label   string
	cursor  int
	choice  chan Choice
}

func Run(label string, choices []Choice) (*Choice, error) {
	result := make(chan Choice, 1)

	if len(choices) == 0 {
		return nil, eris.Wrap(errs.NotFound, "no choices found")
	} else if len(choices) == 1 {
		result <- choices[0]
	} else {
		p := tea.NewProgram(initialModel(label, choices, result))
		if err := p.Start(); err != nil {
			return nil, err
		}
	}

	r := <-result

	return &r, nil
}

type Choice struct {
	Title string
	Value string
}

func initialModel(label string, choices []Choice, choice chan Choice) model {
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
			close(m.choice) // If we're quitting just chose the channel.
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
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
