package termboxUI

import (
	"github.com/nsf/termbox-go"
)

// A button is a simple button field that can be cast to a command.
type Button struct {
	Text   string
	Height int
	Width  int
	Event  UIEvent
	Fg     termbox.Attribute
	Bg     termbox.Attribute
	Active bool
}

// Creates an instance of a Button
func CreateButton(width, height int, text string, fg, bg termbox.Attribute) *Button {
	button := new(Button)
	button.Text = text
	button.Width = width
	button.Height = height
	button.Fg = fg
	button.Bg = bg
	return button
}

func (b *Button) Draw(x, y int) {
	fg := b.Fg
	bg := b.Bg

	if b.Active {
		if b.Bg == termbox.ColorDefault {
			fg = termbox.ColorWhite
		} else {
			fg = b.Bg
		}
		if b.Fg == termbox.ColorDefault {
			bg = termbox.ColorBlack
		} else {
			bg = b.Fg
		}

	}

	textbox := CreateTextBox(b.Width, b.Height, true, false, TextAlignmentDefault, TextAlignmentDefault, fg, bg)
	textbox.AddText(b.Text)
	textbox.Draw(x, y)
}

func (b *Button) HandleKey(key termbox.Key, ch rune, event chan UIEvent) bool {
	switch key {
	case termbox.KeyEnter:
		event <- b.Event
		return true
	default:
		return false
	}
}
