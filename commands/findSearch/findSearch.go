package findSearch

import (
	"fmt"
	"github.com/hjertnes/roam/utils"

	"github.com/charmbracelet/bubbles/textinput"
	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Run() (string ,error){
	result := make(chan string, 1)

	p := tea.NewProgram(initialModel(result))
	err := p.Start()
	utils.ErrorHandler(err)

	if r := <-result; r != "" {
		return r, nil
	}

	return "", nil
}

type errMsg error

type model struct {
	data chan string
	textInput input.Model
	err       error
}

func initialModel(data chan string) model {
	inputModel := input.NewModel()
	inputModel.Placeholder = "My super awesome note"
	inputModel.Focus()
	inputModel.CharLimit = 156
	inputModel.Width = 20

	return model{
		textInput: inputModel,
		err:       nil,
		data: data,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			fallthrough
		case tea.KeyEnter:
			m.data <- m.textInput.Value()
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}


	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Search\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
