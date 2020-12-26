package textinput

import (
	"fmt"
	"github.com/rotisserie/eris"
	"github.com/charmbracelet/bubbles/textinput"
	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(placeholder, label string) (string, error){
	result := make(chan string, 1)

	p := tea.NewProgram(initialModel(result, placeholder, label))
	err := p.Start()

	if err != nil{
		return "", eris.Wrap(err, "textinput failed")
	}

	if r := <-result; r != "" {
		return r, nil
	}

	return "", nil
}

type errMsg error

type model struct {
	label string
	data chan string
	textInput input.Model
	err       error
}

func initialModel(data chan string, placeholder, label string) model {
	inputModel := input.NewModel()
	inputModel.Placeholder = placeholder
	inputModel.Focus()
	inputModel.CharLimit = 156
	inputModel.Width = 20

	return model{
		label: label,
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
		"%s\n\n%s\n\n%s",
		m.label,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
