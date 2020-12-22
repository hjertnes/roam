package findSelect

import (
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

var (
	choices = []string{}
)

type model struct {
	cursor int
	choice chan string
}

func Run(matches []string, write bool, dal *dal.Dal) {
	choices = matches
	result := make(chan string, 1)

	if len(matches)== 0{
		fmt.Println("No matches")
		return
	} else if len(matches) == 1{
		result <- matches[0]
	} else {
		p := tea.NewProgram(model{cursor: 0, choice: result})
		if err := p.Start(); err != nil {
			fmt.Println("Oh no:", err)
			os.Exit(1)
		}
	}


	if r := <-result; r != "" {
		if write {
			editor := utils.GetEditor()
			cmd := exec.Command(editor, r)

			err := cmd.Run()
			utils.ErrorHandler(err)

			err = dal.SetOpened(r)
			utils.ErrorHandler(err)
		} else {
			data, err := ioutil.ReadFile(r)
			utils.ErrorHandler(err)
			metadata := models.Fm{}
			err = frontmatter.Unmarshal(data, &metadata)
			r, _ := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
			)

			out, err := r.Render(fmt.Sprintf("# %s\n%s", metadata.Title, metadata.Content))
			fmt.Print(out)
		}

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