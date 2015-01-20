package termboxUI

//TODO: text box doesn't parse `\n\n` correctly. All text following these characters is ignored.

import (
	"bufio"
	"io"
	"strings"

	"github.com/nsf/termbox-go"
)

//======================================================//
// Useful Functions
//======================================================//

// This fills all of the cells of the terminal within a given rectangle to the specified attributes.
func FillArea(x, y, w, h int, fg, bg termbox.Attribute) {
	for row := 0; row < h; row++ {
		for column := 0; column < w; column++ {
			termbox.SetCell(x+column, y+row, ' ', fg, bg)
		}
	}
	return
}

// Draws a line to the terminal starting with the cell located at 'x' and continuing to the cell at 'w'
// Cells 'x' and 'w' are included.
func DrawHorizontalLine(x, y, w int, fg, bg termbox.Attribute) {
	for i := 0; i <= w; i++ {
		termbox.SetCell(x+i, y, '─', fg, bg)
	}
	return
}

// Draws a line to the terminal starting with the cell located at 'y' and continuing to the cell at 'h'
// Cells 'y' and 'h' are included.
func DrawVerticalLine(x, y, h int, fg, bg termbox.Attribute) {
	for i := 0; i <= h; i++ {
		termbox.SetCell(x, y+i, '│', fg, bg)
	}
	return
}

// Like FillArea, but it also draws a border around the area using the 'fg' attribute as the color.
func DrawRectangle(x, y, h, w int, fg, bg termbox.Attribute) {
	FillArea(x, y, w, h, fg, bg)
	DrawHorizontalLine(x, y, w, fg, bg)    // top
	DrawHorizontalLine(x, h+y, w, fg, bg)  // bottom
	DrawVerticalLine(x, y, h, fg, bg)      // left
	DrawVerticalLine(x+w, y, h, fg, bg)    // right
	termbox.SetCell(x, y, '┌', fg, bg)     // top-left corner
	termbox.SetCell(x+w, y, '┐', fg, bg)   // top-right corner
	termbox.SetCell(x, h+y, '└', fg, bg)   // bottom-left corner
	termbox.SetCell(x+w, h+y, '┘', fg, bg) // bottom-right corner
}

//======================================================//
// Basic Text
//======================================================//

// This is the most basic text drawing function.
// It writes a single line of text to the terminal with the specified settings.
func DrawText(x, y int, line string, fg, bg termbox.Attribute) (int, int) {
	for i, ch := range line {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
	return x + len(line), y
}

// This returns the termbox x coordinate to center the given string within the described area.
// That coordinate value returned should be referenced before drawing the text.
// Note that this doesn't actually draw the text string to the terminal.
func HorizontalCenterString(text string, dimension, offset int) int {
	return (dimension-len(text))/2 + offset
}

//======================================================//
// Text Box
//======================================================//

// Values for passing in text alignment for text boxes.
// TexAlignmentDefault can be used in both vertical and horizontal cases when drawing the text box to the terminal.
const (
	TextAlignmentLeft uint16 = iota
	TextAlignmentRight
	TextAlignmentCenter
	TextAlignmentTop
	TextAlignmentBottom

	TextAlignmentDefault
)

// Basic text box for displaying text in a termbox window.
// HasBorder indicates that the border around the text box should be included when drawing. Note that the borders are drawn within the defined text box's area, effectively losing two columns and two rows of text writing area.
type TextBox struct {
	HasBorder                   bool
	WrapText                    bool
	TextHorizontalJustification uint16
	TextVerticalJustification   uint16
	Width                       int
	Height                      int
	Default_fg                  termbox.Attribute
	Default_bg                  termbox.Attribute

	text        []string
	textHeight  int
	activeIndex int
	scrolling   bool
	reader      io.Reader
}

// This will create a new text box definition.
// If the width or height exceed the dimensions of the termbox, then the screen dimension will be used in place of 'width' or 'height'
func CreateTextBox(width, height int, withBorder, wrapText bool, justification_h, justification_v uint16, fg, bg termbox.Attribute) *TextBox {
	textbox := new(TextBox)
	screenWidth, screenHeight := termbox.Size()

	if width == -1 || width > screenWidth {
		textbox.Width = screenWidth
	} else {
		textbox.Width = width
	}

	if height == -1 {
		textbox.Height = screenHeight
	} else {
		textbox.Height = height
	}

	textbox.TextHorizontalJustification = justification_h
	textbox.TextVerticalJustification = justification_v

	textbox.Default_fg = fg
	textbox.Default_bg = bg

	textbox.HasBorder = withBorder
	textbox.WrapText = wrapText

	newHeight := textbox.Height
	if textbox.HasBorder && textbox.Height > 2 {
		newHeight -= 2
	}
	textbox.text = make([]string, newHeight)
	textbox.textHeight = 0
	textbox.activeIndex = 0
	textbox.scrolling = true
	textbox.reader = nil

	return textbox
}

// This lets a text box accept a reader instead of an explicit string.
// The assumption is that the type of data from the read source is always 'string', at least for now...
func (tb *TextBox) AddTextFrom(strReader io.Reader) error {
	tb.reader = strReader
	return nil
}

// This adds a single line of text to the text box.
// The '\n' rune is translated to a new line and the '\t' rune is treated as four spaces.
func (tb *TextBox) AddText(text string) {
	height := tb.Height
	width := tb.Width

	if tb.HasBorder {
		height -= 2
		width -= 2
	}

	var lines []string
	linesHeight := 0
	strArray := strings.Split(text, "\n")
	for _, line := range strArray {

		line = strings.Replace(line, "\t", "    ", -1)

		if len(line) == 0 {
			break
		}

		if tb.WrapText && len(line) > width {
			for len(line) != 0 {
				var newLine = ""

				if len(line) < width {
					newLine = line
					line = ""
				} else {
					newLine = line[:width-1]
					line = line[width:]
				}

				if !tb.scrolling && linesHeight+tb.textHeight <= height {
					lines = append(lines, newLine)
					linesHeight++
				} else {
					break
				}
			}
		} else {
			if !tb.scrolling && linesHeight+tb.textHeight <= height {
				lines = append(lines, line)
				linesHeight++
			} else {
				lines = append(lines, line)
				linesHeight++
			}
		}
	}

	if linesHeight > 0 {
		temp := make([]string, linesHeight+tb.textHeight)
		copy(temp, tb.text)
		copy(temp[tb.textHeight:], lines)
		tb.text = temp

		tb.textHeight += linesHeight
	}
}

// This will write the text box to the terminal. 'x' and 'y' are the upper-left coordinates from which the box will be drawn.
// The cell at that location is included when drawing.
// If the number of lines of the text box after wrapping is applied is larger than the height of the box, scrolling is automatically applied.
func (tb *TextBox) Draw(x, y int) {

	if tb.reader != nil {
		reader := bufio.NewReader(tb.reader)
		str, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if len(str) > 0 {
			tb.AddText(str)
		}
	}

	width := tb.Width
	height := tb.Height

	if tb.HasBorder {
		DrawRectangle(x, y, height, width, tb.Default_fg, tb.Default_bg)
		width -= 2
		height -= 2
		x++
		y++
	} else {
		FillArea(x, y, width, height, tb.Default_fg, tb.Default_bg)
	}

	for i := 0; i <= height; i++ {
		if i+tb.activeIndex > tb.textHeight {
			break
		}

		if tb.activeIndex+i == len(tb.text) {
			break
		}
		var x_coord, y_coord int

		line := tb.text[tb.activeIndex+i]

		switch tb.TextHorizontalJustification {
		case TextAlignmentCenter:
			x_coord = HorizontalCenterString(line, width, x)
		case TextAlignmentRight:
			x_coord = (x + width) - len(line) - 1
		default:
			x_coord = x
		}

		switch tb.TextVerticalJustification {
		case TextAlignmentCenter:
			y_coord = y + (height / 2) + i
		case TextAlignmentBottom:
			y_coord = (y + height) - tb.textHeight + i
		default:
			y_coord = y + i
		}

		DrawText(x_coord, y_coord, line, tb.Default_fg, tb.Default_bg)
	}
}

// Handle the termbox key or character input.
// The up and down keys will scroll the text within the text box area.
// Any other input is ignored by text box.
func (tb *TextBox) HandleKey(key termbox.Key, ch rune, results chan UIEvent) bool {
	eventConsumed := true

	switch key {
	case termbox.KeyArrowUp:
		tb.activeIndex--
		if tb.activeIndex < 0 {
			tb.activeIndex = 0
		}
	case termbox.KeyArrowDown:
		if !(tb.activeIndex+tb.Height >= tb.textHeight+1) {
			tb.activeIndex++
		}
	default:
		eventConsumed = false
	}

	return eventConsumed
}
