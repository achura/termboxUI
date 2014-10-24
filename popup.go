package termboxUI

import (
	"github.com/nsf/termbox-go"
)

// Values associated with the Popup field.
const (
	// These are the different types of popups.
	// All of these types allow for a title and content text.
	InputPopup   uint16 = iota // editbox for user text input.
	OKPopup                    // a single 'ok' button. Press any key to hide the popup.
	YesNoPopup                 // 'yes' and 'no' buttons
	DefaultPopup               // default. Does not handle any user input and is always drawn when it is part of a UI.

	// These indicate on which side, top or bottom, the popup will be drawn. Default will draw the popup centered on the terminal window.
	PopupBottom uint16 = iota
	PopupTop
	PopupDefault
)

// A popup is a simple text box that can be justified to the top, bottom or centered
// It can be just a title or a title with a line break and limited text
type Popup struct {
	Title    string
	Content  string
	Position uint16
	Button1  Button
	Button2  Button
	Type     uint16
	Width    int
	Height   int
	Fg       termbox.Attribute
	Bg       termbox.Attribute
}

func CreatePopup(title, content string, position /*, pType*/ uint16, height, width int, fg, bg termbox.Attribute) *Popup {
	popup := new(Popup)

	if position == PopupTop || position == PopupBottom {
		popup.Position = position
	} else {
		popup.Position = PopupDefault
	}

	screenWidth, screenHeight := termbox.Size()

	popup.Width = width
	if width == -1 {
		popup.Width = screenWidth
	}

	popup.Height = height
	if height == -1 {
		popup.Height = screenHeight
	}

	popup.Title = title
	popup.Content = content
	popup.Fg = fg
	popup.Bg = bg

	return popup
}

// This will draw the popup to the terminal.
// Note that because popups are static fields, the x and y input values are ignored when drawing.
// They are included as input options so that the Popup struct is a DrawHandler interface.
func (pu *Popup) Draw(x, y int) {
	textBox := CreateTextBox(pu.Width, pu.Height, true, true, TextAlignmentCenter, TextAlignmentDefault, pu.Fg, pu.Bg)

	screenWidth, screenHeight := termbox.Size()
	x = (screenWidth - pu.Width) / 2
	y = (screenHeight - pu.Height) / 2

	if pu.Position == PopupTop {
		y = -1
	} else {
		y = screenHeight - pu.Height + 1
	}

	textBox.AddText(pu.Title)

	if len(pu.Content) > 0 {
		bar := ""
		for i := 0; i < pu.Width-4; i++ {
			bar += "─"
		}

		textBox.AddText(bar)
		textBox.AddText(pu.Content)
	}

	textBox.Draw(x, y)
}

// Currently, the popup does not take any input
func (pu *Popup) HandleKey(key termbox.Key, ch rune, event chan UIEvent) bool { return false }
