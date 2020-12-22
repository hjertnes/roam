package findSelect

import (
	"fmt"
	"github.com/hjertnes/roam/utils"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	choices = []string{}
)

type model struct {
	cursor int
	choice chan string
}

func Run(matches []string) {
	choices = matches
	result := make(chan string, 1)
	p := tea.NewProgram(model{cursor: 0, choice: result})
	if err := p.Start(); err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	if r := <-result; r != "" {
		editor := utils.GetEditor()
		cmd := exec.Command(editor, r)

		err := cmd.Start()
		utils.ErrorHandler(err)
	}
}

func initialModel(choice chan string) model {
	return model{cursor: 0, choice: choice}
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
			m.choice <- choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}

	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("Select file to open\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}