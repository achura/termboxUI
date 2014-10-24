package termboxUI

import (
	"github.com/nsf/termbox-go"
)

//==========================//
//       Menu Option        //
//==========================//

// MenuOption represents a single menu item.
// The Title value will be the text displayed for that option and must be populated.
// HelpText is optional text. This will show in a popup when 'F1' is pressed if menu help text is supported.
type MenuOption struct {
	Title    string
	HelpText string
	Command  func() UIEvent
}

// Executes the function represented by the menu option and returns the result on an event channel.
func (mo *MenuOption) ExecuteCommand(ev chan UIEvent) {
	ev <- mo.Command()
}

//==========================//
//           Menu           //
//==========================//

// Menus can be displayed as either a grid or in a list.
type MenuMode int

// Values associated with the Menu field.
const (
	MenuList MenuMode = iota
	MenuGrid

	// If this is used in place of an index, the MenuOption will be appended to the list.
	MenuInsertLast int = -1
)

// A Menu is a fully-featured menu for the termbox-go platform!
// This consists of an array of options, each one capable executing its own command to handle user interaction
// A user can either use the arrow keys or a number to highlight a menu option. Use the 'enter' or 'return' key to select that option and execute its command.
type Menu struct {
	Width       int
	Height      int
	Header      string
	Mode        MenuMode
	DrawHelpBox bool
	Fg          termbox.Attribute
	Bg          termbox.Attribute

	Options []MenuOption

	activeIndex int
	menuTop     int
	menuBottom  int
}

// This creates an instance of a new Menu.
// If drawHelpBox is true then the F1 key will display the description of the menu option using a pop up at the bottom of the screen.
func CreateMenu(width, height int, header string, mode MenuMode, drawHelpBox bool, fg, bg termbox.Attribute) *Menu {
	options := make([]MenuOption, 0)
	return &Menu{width, height, header, mode, drawHelpBox, fg, bg, options, 0, 0, height}
}

// this adds a new menu option
// Either use MenuInsertLast to add the new option to the bottom of the options list or specify the index of the new menu option.
func (m *Menu) InsertMenuOption(index int, newOption MenuOption) {
	var newOptions = make([]MenuOption, len(m.Options)+1)

	if 0 <= index && index < len(m.Options) {
		copy(newOptions[:index], m.Options[:index+1])
		copy(newOptions[:index+1], m.Options[index:])
		newOptions[index] = newOption
		m.Options = newOptions
	} else {
		copy(newOptions, m.Options)
		newOptions[len(m.Options)] = newOption
		m.Options = newOptions
	}
}

// Remove the menu option at the specified index from the menu
func (m *Menu) RemoveMenuOption(index int) {
	copy(m.Options[index:], m.Options[index+1:])
	m.Options = m.Options[:len(m.Options)-1]
}

// Replace the menu option at the specified index with a new option
// Does nothing if the index is not within the current range of indices.
// Todo: only use valid options
func (m *Menu) ReplaceMenuOption(index int, newOption MenuOption) {
	if 0 <= index && index < len(m.Options) {
		m.Options[index] = newOption
	}
}

// Draws the menu to the terminal at the specified indices.
func (m *Menu) Draw(x, y int) {
	cols := 1

	//Draw the menu Title
	if len(m.Header) > 0 {
		titleBox := CreateTextBox(m.Width, 1, false, false, TextAlignmentCenter, TextAlignmentDefault, m.Fg, m.Bg)
		titleBox.AddText(m.Header)
		titleBox.Draw(x, y)
		DrawHorizontalLine(x, y+1, m.Width, m.Fg, m.Bg)
		y += 3
	}

	if len(m.Options) < m.Height {
		m.menuBottom = len(m.Options)
	}

	rows := m.menuBottom

	if m.Mode == MenuGrid {
		cols = 2
		rows = m.menuBottom/cols + 1
	}

	table := CreateTable(m.Width, m.menuBottom, cols, rows, nil, nil, false, true, m.Fg, m.Bg)
	for c := 0; c < cols; c++ {
		for r := m.menuTop; r < rows; r++ {
			index := getIndexFromCoordinates(rows, c, r)

			if index >= m.menuBottom {
				break
			}

			table.SetCell(c, r, m.Options[index].Title)
		}
	}

	table.ActiveColumn, table.ActiveRow = getCoordinatesFromIndex(rows, m.activeIndex)

	y -= m.menuTop
	table.Draw(x, y)

	if m.DrawHelpBox {
		drawHelpBox(m.Options[m.activeIndex].HelpText, m.Fg, m.Bg)
	}
}

// Handles input termbox key or character.
// The arrow keys will change the active or highlighted menu option.
// A number key will select the option at the specified index.
// If help text is enabled, 'F1' will toggle the help text as a popup from the bottom of the terminal.
// Any other user input is ignored.
func (m *Menu) HandleKey(key termbox.Key, ch rune, results chan UIEvent) (eventConsumed bool) {
	eventConsumed = true

	switch key {
	case termbox.KeyArrowUp:
		m.activeIndex--
		if m.activeIndex < 0 {
			m.activeIndex = 0
		}
		if m.activeIndex == m.menuTop && m.activeIndex > 0 {
			m.menuTop--
			m.menuBottom--
		}
	case termbox.KeyArrowDown:
		m.activeIndex++
		if m.activeIndex >= len(m.Options) {
			m.activeIndex = len(m.Options) - 1
		}
		if m.activeIndex == m.menuBottom-1 && m.activeIndex < len(m.Options)-1 {
			m.menuTop++
			m.menuBottom++
		}
	case termbox.KeyArrowLeft:
		if m.Mode == MenuGrid && m.activeIndex >= len(m.Options)/2 {
			m.activeIndex -= len(m.Options)/2 + 1
		}
	case termbox.KeyArrowRight:
		if m.Mode == MenuGrid && m.activeIndex <= len(m.Options)/2 && m.activeIndex+2 < len(m.Options) {
			m.activeIndex += len(m.Options)/2 + 1
		}
	case termbox.KeyEnter:
		go m.Options[m.activeIndex].ExecuteCommand(results)
	case termbox.KeyF1:
		m.DrawHelpBox = !m.DrawHelpBox
	default:
		//If it is a number, set that as the active index
		if ch != 0 {
			for index, char := range "123456789" {
				if ch == char {
					m.activeIndex = index
					break
				}
			}
		} else {
			eventConsumed = false
		}
	}
	return
}

//==========================//
//        Utilities         //
//==========================//

// Convert from the termbox window coordinates to the index of a menu option at that coordinate.
func getIndexFromCoordinates(rows, col, row int) int {
	return rows*col + row
}

// Convert from the menu option at the specified index to the termbox cell coordinates.
func getCoordinatesFromIndex(rows, index int) (col, row int) {
	row = (index) % rows
	col = (index - row) / rows
	return
}

// Displays the popup with the menu help text.
func drawHelpBox(text string, fg, bg termbox.Attribute) {
	popup := CreatePopup("ABOUT", text, PopupBottom, 6, -1, fg, bg)
	popup.Draw(0, 0)
	return
}
