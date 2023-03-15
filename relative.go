package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	command       string
	currentFolder string
	commandOut    []string
}

func initialModel() model {
	return model{
		command:       "",
		currentFolder: "/",
		commandOut:    make([]string, 0),
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
			return m, tea.Quit
		case "enter", " ":
			fmt.Println(m.command)
			m.command = ""
		default:
			m.command += msg.String()
		}
	}
	return m, nil
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
	folderLocation := "Folder location current: " + out.String()
	s += fmt.Sprintf("\n\n%s\n", folderLocation)
	s += "\nPress q to quit.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
