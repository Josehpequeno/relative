package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var primaryColor = tcell.NewRGBColor(25, 36, 50)
var secundaryColor = tcell.NewRGBColor(255, 250, 224)

var app = tview.NewApplication()

var flex = tview.NewFlex()
var pages = tview.NewPages()

func main() {

	tview.Styles.PrimitiveBackgroundColor = primaryColor
	cmd := exec.Command("pwd")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("command stdout: %q", out.String())
	folderLocation := "Folder location current: " + out.String()

	var text = tview.NewTextView().
		SetTextColor(secundaryColor).
		SetText(folderLocation)

	// _box := tview.NewBox().SetBorder(true).SetBorderColor(primaryColor).SetBackgroundColor(primaryColor).SetTitle(" Relative ").SetTitleColor(secundaryColor)
	flex.SetDirection(tview.FlexRow).AddItem(text, 1, 1, false)

	pages.AddPage("Home", flex, true, true).SetBorder(true).SetBorderColor(primaryColor).SetBackgroundColor(primaryColor).SetTitle(" Relative ").SetTitleColor(secundaryColor)

	if err := tview.NewApplication().SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
