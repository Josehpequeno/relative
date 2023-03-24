package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("FFFAE0")).Background(lipgloss.Color("002236"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	// str := fmt.Sprintf("%d. %s", index+1, i)
	str := fmt.Sprintf("%v", i) // items

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "v":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			if !strings.Contains(m.choice, "/") {
				str := "cat " + m.choice
				cmd := exec.Command("bash", "-c", str)
				cmd.Stdin = strings.NewReader("")
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(out.String())
			}
		case "enter", " ":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			// return m, nil
			if strings.Contains(m.choice, "/") {
				os.Chdir(m.choice)

				cmd := exec.Command("ls", "-ap")
				cmd.Stdin = strings.NewReader("")
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}

				items := []list.Item{}
				outArray := strings.Split(out.String(), "\n")
				// fmt.Println(out.String())
				for i := 1; i < len(outArray); i++ {
					if strings.Contains(outArray[i], "/") || strings.Contains(outArray[i], ".") {
						items = append(items, item(outArray[i]))
					}
				}

				const defaultWidth = 20

				l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
				l.Title = "Relative"
				l.SetShowStatusBar(false)
				l.SetFilteringEnabled(false)
				l.Styles.Title = titleStyle
				l.Styles.PaginationStyle = paginationStyle
				l.Styles.HelpStyle = helpStyle
				model := model{list: l}
				m = model
				// return "\n" + model.list.View()
				// model.list.View()
			} else {
				switch choice := strings.Split(m.choice, ".")[1]; choice {
				case "go":
					// return quitTextStyle.Render("is a go file")
					fmt.Println("is a go file")
				case "js":
					fmt.Println("is a javascript file")
				case "ts":
					fmt.Println("is a typescript file")
				case "py":
					fmt.Println("is a python file")
				case "rs":
					fmt.Println("is a rust file")
				default:
					// return quitTextStyle.Render("unrecognized file type")
					fmt.Println("unrecognized file type")
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// if m.choice != "" {
	// return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	// if strings.Contains(m.choice, "/") {
	// 	// fmt.Println("is a directory")
	// 	// str := "cd ./" + m.choice[:len(m.choice)-1]
	// 	// cmd := exec.Command("bash", "-c", str)
	// 	// cmd.Stdin = strings.NewReader("")
	// 	// var out bytes.Buffer
	// 	// cmd.Stdout = &out
	// 	// err := cmd.Run()
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// 	os.Chdir(m.choice)

	// 	cmd := exec.Command("ls", "-ap")
	// 	cmd.Stdin = strings.NewReader("")
	// 	var out bytes.Buffer
	// 	cmd.Stdout = &out
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	items := []list.Item{}
	// 	outArray := strings.Split(out.String(), "\n")
	// 	// fmt.Println(out.String())
	// 	for i := 0; i < len(outArray); i++ {
	// 		if strings.Contains(outArray[i], "/") || strings.Contains(outArray[i], ".") {
	// 			items = append(items, item(outArray[i]))
	// 		}
	// 	}

	// 	const defaultWidth = 20

	// 	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	// 	l.Title = "Relative"
	// 	l.SetShowStatusBar(false)
	// 	l.SetFilteringEnabled(false)
	// 	l.Styles.Title = titleStyle
	// 	l.Styles.PaginationStyle = paginationStyle
	// 	l.Styles.HelpStyle = helpStyle
	// 	model := model{list: l}
	// 	return "\n" + model.list.View()
	// 	// model.list.View()
	// } else {
	// 	switch choice := strings.Split(m.choice, ".")[1]; choice {
	// 	case "go":
	// 		return quitTextStyle.Render("is a go file")
	// 		// fmt.Println("is a go file")
	// 	default:
	// 		return quitTextStyle.Render("unrecognized file type")
	// 		// fmt.Println("unrecognized file type")
	// 	}

	// }
	// }
	if m.quitting {
		return quitTextStyle.Render("bye")
	}
	return "\n" + m.list.View()
}

func main() {
	cmd := exec.Command("ls", "-ap")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}
	outArray := strings.Split(out.String(), "\n")
	for i := 1; i < len(outArray); i++ {
		if strings.Contains(outArray[i], "/") || strings.Contains(outArray[i], ".") {
			items = append(items, item(outArray[i]))
		}
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Relative"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
