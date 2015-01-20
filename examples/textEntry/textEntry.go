package main

import (
	"strings"

	"github.com/C2FO/termboxUI"
	"github.com/nsf/termbox-go"
)

var userText string = "example string"

const ChangeUserText uint16 = iota

func main() {
	if err := termboxUI.StartUI(buildUserInterface); err != nil {
		panic(err)
	}
}

func buildUserInterface() *termboxUI.UI {
	var x, y int
	screenWidth, _ := termbox.Size()

	ui := new(termboxUI.UI)

	// Headline
	title := "Input your message in the box below.\n \nPress `Enter` to display your input all funky and whatnot.\nPress `Esc` to quit."
	headline := termboxUI.CreateTextBox(len(title)+2, 7, false, false, termboxUI.TextAlignmentCenter, termboxUI.TextAlignmentDefault, termbox.ColorDefault, termbox.ColorDefault)
	headline.AddText(title)
	x = (screenWidth - headline.Width) / 2
	y = 1
	ui.AddField(headline, x, y, false)

	// User field
	userField := termboxUI.CreateTextBox(screenWidth-2, 3, false, false, termboxUI.TextAlignmentCenter, termboxUI.TextAlignmentCenter, termbox.ColorDefault, termbox.ColorDefault)
	userField.AddText(funkifyString(userText))
	y = y + 4
	ui.AddField(userField, 1, y, false)

	// Input Box
	inputBox := termboxUI.CreateEditBox(30, userText, ChangeUserText, termbox.ColorDefault, termbox.ColorDefault)
	x = (screenWidth - inputBox.Width) / 2
	y = y + 3
	ui.AddField(inputBox, x, y, true)

	ui.Fg = termbox.ColorDefault
	ui.Bg = termbox.ColorDefault

	// Event Handlers
	ui.CustomEvents = make(map[uint16]func(termboxUI.UIEvent))
	ui.CustomEvents[ChangeUserText] = func(event termboxUI.UIEvent) {
		userText = event.Data.String()
	}

	return ui
}

// This will swap spaces with underscores and swap each character's case.
func funkifyString(str string) string {
	newString := make([]rune, len(str))
	for i, ch := range str {
		if ch == ' ' {
			newString[i] = '_'
		} else if ch == '_' {
			newString[i] = ' '
		} else {
			newString[i] = swapCase(ch)
		}
	}
	return string(newString)
}

func swapCase(ch rune) rune {
	str := string(ch)
	if strings.Contains(strings.ToLower(str), str) {
		return rune(strings.ToUpper(str)[0])
	} else {
		return rune(strings.ToLower(str)[0])
	}
}
