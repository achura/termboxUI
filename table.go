package termboxUI

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// tableCoordinates are the grid coordinates for representing a cell of the table.
// These should not be confused for the termbox.Cell coordinates that are used for drawing to the terminal window.
type tableCoordinates struct {
	X, Y int
}

// This is the representation of a single cell within the table.
type tableCell struct {
	value string
}

// Data within the table is stored represented in rows of cells.
type tableRow []tableCell

// This is a spreadsheet/table for the termbox-go library.
// If ActiveRow and ActiveColumn are both set, the table coordinate they represent will be the only one highlighted. If that location does not lay within the table definitions, only the valid row or column if either with be highlighted.
type Table struct {
	Height       int
	Width        int
	Columns      int
	Rows         int
	Fg           termbox.Attribute
	Bg           termbox.Attribute
	ColumnLabels []string
	RowLabels    []string
	ShowGrid     bool
	ShowNumbers  bool
	ActiveRow    int
	ActiveColumn int

	cells []tableRow
}

// Creates an instance of a new table or spreadsheet.
// If the number of rows exceeds the height of the table, the row count is set to the height.
func CreateTable(width, height, columns, rows int, columnLabels, rowLabels []string, showGrid, showNumbers bool, fg, bg termbox.Attribute) *Table {
	table := new(Table)

	table.Fg = fg
	table.Bg = bg
	table.ShowGrid = showGrid
	table.Rows = rows
	table.Columns = columns
	table.Height = height
	table.Width = width
	table.ShowGrid = showGrid
	table.ShowNumbers = showNumbers

	if height < rows {
		table.Rows = table.Height
	}

	if len(columnLabels) > 0 {
		table.ColumnLabels = make([]string, table.Columns)
		copy(table.ColumnLabels, columnLabels)
	}

	if len(rowLabels) > 0 {
		table.RowLabels = make([]string, table.Rows)
		copy(table.RowLabels, rowLabels)
	}

	table.cells = make([]tableRow, table.Columns)
	for i := range table.cells {
		table.cells[i] = make([]tableCell, table.Rows)
		for _, cell := range table.cells[i] {
			cell.value = ""
		}
	}
	table.ActiveRow = -1
	table.ActiveColumn = -1

	return table
}

// Sets the value of the cell at the specified column and row.
// The return value is 'false' if the column and row coordinates are not within the table parameters.
func (t *Table) SetCell(column, row int, text string) bool {
	if column >= t.Columns || row >= t.Rows || column < 0 || row < 0 {
		return false
	}

	t.cells[column][row].value = text

	return true
}

// Draws the table to the terminal.
// Note that the normal textbox rules for border and dimensions apply to the table.
func (t *Table) Draw(x, y int) {
	number := 0
	cellWidth := (t.Width + 2*len(t.cells)) / t.Columns
	cellHeight := t.Height / t.Rows

	for i, column := range t.cells {
		//Calculate the x-coordinate by making sure that the cells overlap by one character block so that they can share a single line when a grid is active.
		x_coord := x + i*cellWidth - i*2

		for j, row := range column {
			skip := false
			y_coord := y + j*cellHeight

			text := row.value

			if text == "" {
				skip = true
			}

			h_justification := TextAlignmentCenter
			if t.ShowNumbers {
				number++
				h_justification = TextAlignmentLeft
				text = fmt.Sprintf(" %d. %s", number, text)
			}

			fg := t.Fg
			bg := t.Bg

			// Invert the fg and bg colors of any active cell so that it appears highlighted.
			if cellIsActive(t.ActiveColumn, t.ActiveRow, i, j) {
				if t.Bg == termbox.ColorDefault {
					fg = termbox.ColorWhite
				} else {
					fg = t.Bg
				}
				if t.Fg == termbox.ColorDefault {
					bg = termbox.ColorBlack
				} else {
					bg = t.Fg
				}
			}

			if !skip {
				cell := CreateTextBox(cellWidth, cellHeight, t.ShowGrid, false, h_justification, TextAlignmentCenter, fg, bg)
				cell.AddText(text)
				cell.Draw(x_coord, y_coord)
			}
		}
	}
}

// Returns true if the cell at current_col/current_row is active.
func cellIsActive(active_col, active_row, current_col, current_row int) bool {
	col := false
	row := false

	if active_col == -1 && active_row == -1 {
		return false
	}

	if active_col == -1 {
		col = true
	} else if active_col == current_col {
		col = true
	}

	if active_row == -1 {
		row = true
	} else if active_row == current_row {
		row = true
	}
	return col && row
}

// Currently the table does not take any input directly.
func (t *Table) HandleKey(key termbox.Key, ch rune, event chan UIEvent) bool { return false }
