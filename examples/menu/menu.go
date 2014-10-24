package main

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/C2FO/termboxUI"
	"github.com/nsf/termbox-go"
)

const (
	MainMenu uint16 = iota
	FgColorMenu
	BgColorMenu

	MenuChange uint16 = iota
	FgColorChange
	BgColorChange
)

var (
	fgSetting termbox.Attribute
	bgSetting termbox.Attribute

	activeMenu uint16 = MainMenu

	ui termboxUI.UI
)

// This example shows a very basic menu field in action.
// There are three options: change font color, change background color and quit.
// Help text is supported on the options and either Ctrl+C or Esc will exit.
func main() {
	// Initialize the termbox itself.
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// inputEvent is the channel which will signal results from applying user input to the UI
	inputEvent := make(chan termboxUI.UIEvent)
	defer close(inputEvent)

	refresh := true

loop:
	for {
		// Refresh after a menu or color change to redraw the UI
		if refresh {
			ui = buildUserInterface()
			refresh = false
		}
		ui.Draw()

		select {

		// Handle a result from sending input to a UI field.
		case event := <-inputEvent:
			if event.Error != nil {
				panic(event.Error)
			}
			if event.Type == termboxUI.UIResultUint16 {
				handleMenuResult(event.CustomType, event.Data)
			}
			refresh = true

		//Poll for termbox input events on a different channel to prevent blocking the results from a UI interaction
		case ev := <-pollEvent():
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					break loop
				case termbox.KeyCtrlC:
					break loop
				default:
					ui.HandleInput(ev.Key, ev.Ch, inputEvent)
				}
			case termbox.EventResize:
				refresh = true
			}
		}
	}
}

func buildUserInterface() termboxUI.UI {
	var ui termboxUI.UI

	screenWidth, _ := termbox.Size()

	// Display the UI headline at the top by using a single string.
	// This could also be done by using `helpbox.AddText` on each row instead of '+='
	// The single string is used here to display the support for newline characters in the textbox
	var headline = ""
	headline += "          ____                                     "
	headline += "\n        ,'  , `.                                   "
	headline += "\n     ,-+-,.' _ |                                   "
	headline += "\n  ,-+-. ;   , ||              ,---,          ,--,  "
	headline += "\n ,--.'|'   |  ;|          ,-+-. /  |       ,'_ /|  "
	headline += "\n|   |  ,', |  ':  ,---.  ,--.'|'   |  .--. |  | :  "
	headline += "\n|   | /  | |  || /     \\|   |  ,\"' |,'_ /| :  . |  "
	headline += "\n'   | :  | :  |,/    /  |   | /  | ||  ' | |  . .  "
	headline += "\n;   . |  ; |--'.    ' / |   | |  | ||  | ' |  | |  "
	headline += "\n|   : |  | ,   '   ;   /|   | |  |/ :  | : ;  ; |  "
	headline += "\n|   : '  |/    '   |  / |   | |--'  '  :  `--'   \\ "
	headline += "\n;   | |`-'     |   :    |   |/      :  ,      .-./ "
	headline += "\n|   ;/          \\   \\  /'---'        `--`----'     "
	headline += "\n'---'            `----'                            "
	headlineBox := termboxUI.CreateTextBox(51, 14, false, false, termboxUI.TextAlignmentDefault, termboxUI.TextAlignmentDefault, fgSetting, bgSetting)
	headlineBox.AddText(headline)
	ui.AddField(headlineBox, (screenWidth-51)/2, 0, false)

	// Add the menu to the UI
	menu := setMenu(10)
	ui.AddField(menu, 2, 16, true)

	// Set the fg and bg attributes for all fields in the UI
	ui.Fg = fgSetting
	ui.Bg = bgSetting

	return ui
}

func pollEvent() chan termbox.Event {
	event := make(chan termbox.Event)
	go func() {
		event <- termbox.PollEvent()
	}()
	return event
}

func quit() termboxUI.UIEvent {
	termbox.Close()
	os.Exit(0)
	return termboxUI.UIEvent{}
}

func handleMenuResult(customType uint16, data *bytes.Buffer) {
	result := data.Bytes()

	switch customType {
	case MenuChange:
		rdr := bytes.NewReader(result)
		if err := binary.Read(rdr, binary.LittleEndian, &activeMenu); err != nil {
			panic(err)
		}
	case FgColorChange:
		rdr := bytes.NewReader(result)
		if err := binary.Read(rdr, binary.LittleEndian, &fgSetting); err != nil {
			panic(err)
		}
	case BgColorChange:
		rdr := bytes.NewReader(result)
		if err := binary.Read(rdr, binary.LittleEndian, &bgSetting); err != nil {
			panic(err)
		}
	}
}

func setMenu(menuHeight int) (menu *termboxUI.Menu) {
	switch activeMenu {
	case BgColorMenu:
		return getColorMenu(menuHeight, BgColorChange)
	case FgColorMenu:
		return getColorMenu(menuHeight, FgColorChange)
	default:
		return getMainMenu(menuHeight)
	}
}

func getMainMenu(menuHeight int) (menu *termboxUI.Menu) {
	screenWidth, _ := termbox.Size()

	fg_color_option := termboxUI.MenuOption{
		"Font Color",
		"Change the font color.",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, FgColorMenu)

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = MenuChange
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	bg_color_option := termboxUI.MenuOption{
		"Background Color",
		"Change the background color.",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, BgColorMenu)

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = MenuChange
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	exit_option := termboxUI.MenuOption{
		"Quit",
		"Exit the menu example",
		quit,
	}

	menu = termboxUI.CreateMenu(screenWidth-18, menuHeight, "F1 - Toggle help text.", termboxUI.MenuList, false, fgSetting, bgSetting)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, fg_color_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, bg_color_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, exit_option)
	return menu
}

func getColorMenu(menuHeight int, colorChangeType uint16) (menu *termboxUI.Menu) {
	screenWidth, _ := termbox.Size()

	default_option := termboxUI.MenuOption{
		"Default",
		"Use the terminal's default color.",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorDefault))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	black_option := termboxUI.MenuOption{
		"Black",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorBlack))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	white_option := termboxUI.MenuOption{
		"White",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorWhite))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	red_option := termboxUI.MenuOption{
		"Red",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorRed))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	green_option := termboxUI.MenuOption{
		"Green",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorGreen))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	blue_option := termboxUI.MenuOption{
		"Blue",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorBlue))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	yellow_option := termboxUI.MenuOption{
		"Yellow",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorYellow))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	cyan_option := termboxUI.MenuOption{
		"Cyan",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorCyan))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}
	magenta_option := termboxUI.MenuOption{
		"Magenta",
		"Do you seriously need help text here?",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, uint16(termbox.ColorMagenta))

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = colorChangeType
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}

	return_option := termboxUI.MenuOption{
		"Go back",
		"Return to the previous screen",
		func() termboxUI.UIEvent {
			var result = make([]byte, 2)
			binary.LittleEndian.PutUint16(result, MainMenu)

			event := termboxUI.UIEvent{}
			event.Type = termboxUI.UIResultUint16
			event.CustomType = MenuChange
			event.Data = bytes.NewBuffer(result)
			return event
		},
	}

	menu = termboxUI.CreateMenu(screenWidth-18, menuHeight, "Colors", termboxUI.MenuList, false, fgSetting, bgSetting)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, default_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, black_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, white_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, red_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, green_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, blue_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, yellow_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, cyan_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, magenta_option)
	menu.InsertMenuOption(termboxUI.MenuInsertLast, return_option)
	return menu
}
