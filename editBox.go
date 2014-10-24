package termboxUI

import (
	"bytes"

	"github.com/nsf/termbox-go"
)

//============================//
//          Edit Box          //
//----------------------------//

// An edit box or input box is a field that allows the user to input text.
// A custom type can be set to help indicate the nature of the text being input.
// For example, an input box could be for a first name or last name.
type EditBox struct {
	Width       int
	Height      int
	Value       []rune
	Fg          termbox.Attribute
	Bg          termbox.Attribute
	CursorIndex int
	CustomType  uint16
}

// Creates a new instance of an edit box.
// When width is -1, the text box will be the width of the terminal window.
func CreateEditBox(width int, value string, customMessageCode uint16, fg, bg termbox.Attribute) *EditBox {
	editBox := new(EditBox)

	screenWidth, _ := termbox.Size()

	editBox.Width = width
	if editBox.Width == -1 {
		editBox.Width = screenWidth
	}

	editBox.Fg = fg
	editBox.Bg = bg

	editBox.CustomType = customMessageCode

	if len(value) > 0 {
		editBox.Value = make([]rune, len(value))
		copy(editBox.Value, []rune(value))
	} else {
		editBox.Value = make([]rune, 0)
	}
	editBox.CursorIndex = 0

	return editBox
}

func (eb *EditBox) Draw(x, y int) {
	// Convert the input editable string to an array of runes so characters can be added, modified or removed from the edit text.
	if len(eb.Value) > eb.Width {
		nchrs := make([]rune, eb.Width)
		copy(nchrs, eb.Value[len(eb.Value)-eb.Width-1:])
		eb.Value = nchrs
	}
	displayString := string(eb.Value)

	textbox := CreateTextBox(eb.Width, 4, false, false, TextAlignmentDefault, TextAlignmentCenter, eb.Fg, eb.Bg)
	textbox.AddText("/> " + displayString)
	textbox.Draw(x, y)

	x_coord := x + eb.CursorIndex + 3
	termbox.SetCursor(x_coord, y+2)

	return
}

// Handles a termbox key or character input
// 'Enter' or 'Return' will end the string editing and signal that editing is complete.
// 'Backspace' removes the character before the currently selected character.
// 'Delete' removes the currently selected character.
// Left and right arrow keys will move the cursor along the edit string.
// 'Tab' inserts four spaces to the run array.
// 'Space' inserts a single space.
// Otherwise the character input is added to the string.
func (eb *EditBox) HandleKey(key termbox.Key, ch rune, ev chan UIEvent) (eventConsumed bool) {
	eventConsumed = true

	switch key {
	case termbox.KeyEnter:
		// Send along the input
		event := UIEvent{}
		event.Type = UIResultString
		event.CustomType = eb.CustomType
		event.Data = bytes.NewBufferString(string(eb.Value))
		ev <- event

		//Clear the edit buffer
		eb.Value = make([]rune, 0)
		eb.CursorIndex = 0

	case termbox.KeyBackspace2:
		startLength := len(eb.Value)
		eb.Value = removeCharacter(eb.Value, eb.CursorIndex-1)
		if startLength > len(eb.Value) {
			eb.CursorIndex = setCursor(eb.CursorIndex, eb.CursorIndex-1, len(eb.Value))
		}
	case termbox.KeyDelete:
		eb.Value = removeCharacter(eb.Value, eb.CursorIndex)
	case termbox.KeyArrowRight:
		eb.CursorIndex = setCursor(eb.CursorIndex, eb.CursorIndex+1, len(eb.Value))
	case termbox.KeyArrowLeft:
		eb.CursorIndex = setCursor(eb.CursorIndex, eb.CursorIndex-1, len(eb.Value))
	case termbox.KeyTab:
		startLength := len(eb.Value)
		eb.Value = insertCharacter(eb.Value, ' ', eb.CursorIndex)
		eb.Value = insertCharacter(eb.Value, ' ', eb.CursorIndex)
		eb.Value = insertCharacter(eb.Value, ' ', eb.CursorIndex)
		eb.Value = insertCharacter(eb.Value, ' ', eb.CursorIndex)
		if startLength < len(eb.Value) {
			eb.CursorIndex = setCursor(eb.CursorIndex, eb.CursorIndex+4, len(eb.Value))
		}
	case termbox.KeySpace:
		startLength := len(eb.Value)
		eb.Value = insertCharacter(eb.Value, ' ', eb.CursorIndex)
		if startLength < len(eb.Value) {
			eb.CursorIndex = setCursor(eb.CursorIndex, eb.CursorIndex+1, len(eb.Value))
		}
	default:
		if ch != 0 {
			startLength := len(eb.Value)
			eb.Value = insertCharacter(eb.Value, ch, eb.CursorIndex)
			if startLength < len(eb.Value) {
				eb.CursorIndex = setCursor(eb.CursorIndex, eb.CursorIndex+1, len(eb.Value))
			}
		} else {
			eventConsumed = false
		}
	}
	return
}

//============================//
//         Utilities          //
//----------------------------//

// Insert a character into the rune array at a specified index.
// If the index is greater than the length of the edit rune array, it will be appended to the end of the array.
func insertCharacter(dst []rune, ch rune, index int) []rune {
	slice := make([]rune, len(dst)+1)

	switch {
	case index == len(dst):
		//add the character to the end of the array
		copy(slice, dst)
		slice[len(dst)] = ch
		return slice
	case len(dst) > index:
		// insert the character into the array
		copy(slice, dst[:index])
		slice[index] = ch
		copy(slice[index+1:], dst[index:])
		return slice
	default:
		return dst
	}
}

// Remove the character from the rune array that is specifid by 'index.' if it is a valide character within the edit string.
func removeCharacter(dst []rune, index int) []rune {
	if index >= 0 && index < len(dst) {
		slice := make([]rune, len(dst)-1)
		copy(slice, dst[:index])
		copy(slice[index:], dst[index+1:])
		return slice
	}
	return dst
}

// Determine the index of the active/highlighted character in the edit string.
func setCursor(from, to, inputBoxLength int) (newIndex int) {
	newIndex = from

	if from < to && from >= 0 ||
		from > to && to >= 0 {
		newIndex = to
	}

	return
}
