// Package textinput is a wrapper around bubbletea textinput to make it easy to get input from user
package textinput

import (
	"fmt"
	"os"

	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rotisserie/eris"
)

// Run displays the input and returns the result.
func Run(placeholder, label string) (string, error) {
	result := make(chan string, 1)

	p := tea.NewProgram(initialModel(result, placeholder, label))

	err := p.Start()
	if err != nil {
		return "", eris.Wrap(err, "textinput failed")
	}

	if r := <-result; r != "" {
		return r, nil
	}

	return "", nil
}

type errMsg error

type model struct {
	label     string
	data      chan string
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
		label:     label,
		textInput: inputModel,
		err:       nil,
		data:      data,
	}
}

func (m model) Init() tea.Cmd {
	return input.Blink
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
			os.Exit(0)
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
