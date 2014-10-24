// Package termboxUI is a ui library built for termbox-go!
// The fields and tools contained herein help to streamline graphics drawing and user input handling from a terminal.
package termboxUI

import (
	"bytes"

	"github.com/nsf/termbox-go"
)

//==========================//
//         UI Event         //
//==========================//

// Any data reported from a user interaction with a field is sent using a buffer,
// so ResultType is used so that the resulting data can be converted to something usable.
type ResultType uint16

// These are the possible default types returned from a user interaction.
const (
	UIResultBool ResultType = iota
	UIResultByte
	UIResultComplex128
	UIResultComplex64
	UIResultError
	UIResultFloat64
	UIResultFloat32
	UIResultInt
	UIResultInt8
	UIResultInt16
	UIResultInt32
	UIResultInt64
	UIResultRune
	UIResultString
	UIResultUint
	UIResultUint8
	UIResultUint16
	UIResultUint32
	UIResultUint64
	UIResultUintptr

	UIResultJSON
	UIResultXML

	UIResultMap
	UIResultSlice

	UIResultNone
)

// This represents the results from a user interaction with a field.
// The CustomType is a developer-defined type which helps when determining between different
type UIEvent struct {
	Error      error
	Type       ResultType
	CustomType uint16
	Data       *bytes.Buffer
}

// DrawHandler is the basic interface that all ui fields must implement.
// A UI calls Draw to display the DrawHandler at the given x, y location on the terminal.
// A UI sends termbox input data to the DrawHandler through HandleKey to process any user input on a field. It returns 'true' if the message was used.
type DrawHandler interface {
	Draw(x, y int)
	HandleKey(key termbox.Key, ch rune, event chan UIEvent) bool
}

//==========================//
//            UI            //
//==========================//

// All ui fields should adhere to this interface.
type Field struct {
	X        int
	Y        int
	Element  DrawHandler
	HasFocus bool
}

// This is the definition of all of the fields in the current termbox GUI.
type UI struct {
	Fg     termbox.Attribute
	Bg     termbox.Attribute
	fields []Field
}

// AddField adds a new ui field to the defined UI
// The field with draw starting at the specified termbox coordinates
// hasFocus will give the input handling priority to the new field.
func (ui *UI) AddField(element DrawHandler, x, y int, hasFocus bool) {
	var newFields = make([]Field, len(ui.fields)+1)
	var field = Field{x, y, element, hasFocus}
	copy(newFields[:], ui.fields[:])
	newFields[len(ui.fields)] = field
	ui.fields = newFields
	return
}

// Draw clears the terminal and then calls the Draw method for all of its fields at their set locations.
func (ui *UI) Draw() {
	termbox.Clear(ui.Fg, ui.Bg)
	for _, field := range ui.fields {
		field.Element.Draw(field.X, field.Y)
	}
	termbox.Flush()
	return
}

// Send the termbox key and character input to the UI's fields.
// As soon as the event is consumed by a field, this returns. This way only one field can handle that input at a time.
func (ui *UI) HandleInput(key termbox.Key, ch rune, event chan UIEvent) (eventConsumed bool) {
	eventConsumed = false

inputLoop:
	for _, field := range ui.fields {
		if field.HasFocus {
			eventConsumed = field.Element.HandleKey(key, ch, event)
			break inputLoop
		}
	}

	return
}
