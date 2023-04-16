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
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14
const useHighPerformanceRenderer = false

var (
	titleStyle         = lipgloss.NewStyle().MarginLeft(2)
	itemStyle          = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle  = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("FFFAE0")).Background(lipgloss.Color("002236"))
	paginationStyle    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle          = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle      = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	titleStyleViewport = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
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
	content  string
	ready    bool
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		teaCmd  tea.Cmd
		teaCmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		headerHeight := lipgloss.Height(m.headerView(m.choice))
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			// m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			teaCmds = append(teaCmds, viewport.Sync(m.viewport))
		}
		// return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			if !m.ready {
				m.quitting = true
				return m, tea.Quit
			} else {
				m.ready = false
				// m.content = ""
				return m, nil
			}
		case "v":
			m.ready = true
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			if !strings.Contains(m.choice, "/") {
				// 		// str := "cat " + m.choice
				// cmd := exec.Command("pwd")
				// cmd.Stdin = strings.NewReader("")
				// var out bytes.Buffer
				// cmd.Stdout = &out
				// err := cmd.Run()
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// fmt.Println(strings.Trim(out.String(), "\n") + "/" + m.choice)
				// filePath := strings.Trim(out.String(), "\n") + "/" + m.choice
				content, err := os.ReadFile(m.choice)
				if err != nil {
					fmt.Println("could not load file:", err)
					// os.Exit(1)
					m.content = string("could not load file:")
				}
				m.content = string(content)
				// fmt.Println(m.content)
				// 		content, err := os.ReadFile(m.choice)
				// 		if err != nil {
				// 			fmt.Println("could not load file:", err)
				// 			os.Exit(1)
				// 		}

				// 		p := tea.NewProgram(
				// 			model{content: string(content)},
				// 			tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
				// 			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
				// 		)

				// 		if _, err := p.Run(); err != nil {
				// 			fmt.Println("could not run program:", err)
				// 			os.Exit(1)
				// 		}
			} else if m.choice != "../" {
				os.Chdir(m.choice)

				cmd := exec.Command("ls", "-ap")
				cmd.Stdin = strings.NewReader("")
				var out bytes.Buffer
				cmd.Stdout = &out
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				// fmt.Println(out.String())
				m.content = out.String()

				os.Chdir("../")
			} else {
				m.ready = false
			}
			m.viewport.SetContent(m.content)
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
				// model := model{list: l}
				m.list = l
				// return "\n" + model.list.View()
				// model.list.View()
			} else {
				switch choice := strings.Split(m.choice, ".")[1]; choice {
				case "go":
					// return quitTextStyle.Render("is a go file")
					m.ready = true
					str := "go run " + m.choice
					cmd := exec.Command("bash", "-c", str)
					cmd.Stdin = strings.NewReader("")
					var out bytes.Buffer
					cmd.Stdout = &out
					err := cmd.Run()
					if err != nil {
						log.Fatal("error executing go file: ", err)
					}
					m.content = out.String()
					m.viewport.SetContent(m.content)
				case "js":
					m.ready = true
					str := "node " + m.choice
					cmd := exec.Command("bash", "-c", str)
					cmd.Stdin = strings.NewReader("")
					var out bytes.Buffer
					cmd.Stdout = &out
					err := cmd.Run()
					if err != nil {
						log.Fatal("error executing javascript file: ", err)
					}
					m.content = out.String()
					m.viewport.SetContent(m.content)
				case "ts":
					m.ready = true
					str := "node " + m.choice
					cmd := exec.Command("bash", "-c", str)
					cmd.Stdin = strings.NewReader("")
					var out bytes.Buffer
					cmd.Stdout = &out
					err := cmd.Run()
					if err != nil {
						log.Fatal("error executing typescript file: ", err)
					}
					m.content = out.String()
					m.viewport.SetContent(m.content)
				case "py":
					m.ready = true
					str := "python " + m.choice
					cmd := exec.Command("bash", "-c", str)
					cmd.Stdin = strings.NewReader("")
					var out bytes.Buffer
					cmd.Stdout = &out
					err := cmd.Run()
					if err != nil {
						log.Fatal("error executing python file: ", err)
					}
					m.content = out.String()
					m.viewport.SetContent(m.content)
				case "rs":
					m.ready = true
					str := "rustc " + m.choice
					fmt.Println("str", str)
					firstCmd := exec.Command("bash", "-c", str)

					firstCmd.Stdin = strings.NewReader("")
					var out bytes.Buffer
					firstCmd.Stdout = &out
					err := firstCmd.Run()
					if err != nil {
						log.Fatal("error executing rust file: ", err)
					}
					str = "./" + strings.Split(m.choice, ".")[0]
					secondCmd := exec.Command("bash", "-c", str)
					secondCmd.Stdin = strings.NewReader("")
					// var out bytes.Buffer
					secondCmd.Stdout = &out
					err = secondCmd.Run()
					if err != nil {
						log.Fatal("error executing rust file: ", err)
					}
					fmt.Println("str", str)
					// secondCmd = exec.Command("bash", "-c", str)
					m.content = out.String()
					m.viewport.SetContent(m.content)
					str = "rm ./" + strings.Split(m.choice, ".")[0]
					fmt.Println("str", str)
					tertiaryCmd := exec.Command("bash", "-c", str)
					tertiaryCmd.Stdin = strings.NewReader("")
					// var out bytes.Buffer
					tertiaryCmd.Stdout = &out
					err = tertiaryCmd.Run()
					if err != nil {
						log.Fatal("error executing rust file: ", err)
					}
					fmt.Println(out.String())
				default:
					// return quitTextStyle.Render("unrecognized file type")
					fmt.Println("unrecognized file type")
				}
			}
		}
	}
	if m.ready {
		m.viewport, teaCmd = m.viewport.Update(msg)
		teaCmds = append(teaCmds, teaCmd)

		return m, tea.Batch(teaCmds...)
	} else {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
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
	if m.ready {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(m.choice), m.viewport.View(), m.footerView())
	}
	if m.quitting {
		return quitTextStyle.Render("bye")
	}
	return "\n" + m.list.View()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) headerView(name string) string {
	title := titleStyleViewport.Render(name)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
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

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
