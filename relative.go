package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type model struct {
	currentFolder string
	commandOut    []string
	command       textinput.Model
	commandString string
	err           error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "ls"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		command:       ti,
		commandString: "",
		currentFolder: "/",
		commandOut:    make([]string, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func ExecCommand(command string) {
	if command != "" {
		println("command", command)
		cmd := exec.Command(command)
		cmd.Stdin = strings.NewReader("")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println(command, "command not found")
		} else {
			fmt.Println("command out:\n", out.String())
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlQ:
			return m, tea.Quit
		case tea.KeyEnter:
			ExecCommand(m.commandString)
		case tea.KeyBackspace:
			last := len(m.commandString) - 1
			if last >= 0 {
				m.commandString = m.commandString[:last] // remove last char in commandString
			}
			if last == 0 {
				m.commandString = ""
			}
		default:
			m.commandString += msg.String()
			// fmt.Println("> ", m.commandString)
		}
	case errMsg:
		m.err = msg
		return m, nil
		// default:
		// 	m.command += msg.String()
	}

	m.command, cmd = m.command.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := "Relative"
	cmd := exec.Command("pwd")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	s += "\nPress Ctrl+q to quit.\n"
	folderLocation := "Folder location current: " + out.String()
	s += fmt.Sprintf("\n\n%s\n%s", folderLocation, m.command.View())

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
