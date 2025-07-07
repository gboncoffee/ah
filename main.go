package main

import "github.com/gboncoffee/ah/ui"

const scratchString = `Welcome to the ḥăzā' text editor.
This is a scratch buffer you can use for text that will not be saved.

This is the 79th column --->                                                   |
This is the 80th column --->                                                    |

This is a very long line to test wether the text wrapping functionality of the renderer is properly working or not.
`

var UI *ui.Ui

func main() {
	var e Haza

	UI = ui.New()
	e.InitColors()
	e.InitMinibuffer()

	buf := e.NewBuffer(scratchString, "/scratch")

	UI.Start(func(s *ui.State) {
		s.Editor = buf.Editors[0]
		s.Minibuffer = e.Minibuffer.Editor
	})
}
