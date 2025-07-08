package main

import (
	"flag"
	"fmt"

	"github.com/gboncoffee/ah/ui"
)

const debugScratchString = `Welcome to the ḥăzā' text editor.
This is a scratch buffer you can use for text that will not be saved.

This is the 79th column --->                                                   |
This is the 80th column --->                                                    |

This is a very long line to test wether the text wrapping functionality of the renderer is properly working or not.

This is a tab:	it should be 8th-column aligned! And	be	sure	that	lots	of	them	don't	break	proper	wrapping.

	This line is indented with a tab.
`

const scratchString = `Welcome to the ḥăzā' text editor.
This is a scratch buffer you can use for text that will not be saved.


`

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "activates the debug mode")
	flag.Parse()

	E.Ui = ui.New()
	E.InitColors()
	E.InitMinibuffer()

	E.Ui.MessageStyle(E.Colors.Message)
	E.Ui.WarningStyle(E.Colors.Warning)

	E.Ui.EventHandler(func(ev ui.Event) {
		Event(&E, ev)
	})

	E.Ui.Start(func(s *ui.State) {
		var scratch *FileBuffer
		if debug {
			scratch = E.NewBuffer(debugScratchString, "//scratch")
		} else {
			scratch = E.NewBuffer(scratchString, "//scratch")
			// TODO: Proper interface for moving the cursor.
			scratch.editors[0].cursors[0].Begin = 109
			scratch.editors[0].cursors[0].End = 110
		}
		bufToOpen := scratch

		erroedOpening := false
		for _, file := range flag.Args() {
			newBuffer, err := E.Open(file)
			if err != nil {
				erroedOpening = true
				E.LogError(err.Error())
			} else {
				E.Log(fmt.Sprintf("opened file %v", newBuffer.Name()))
				if bufToOpen == scratch {
					bufToOpen = newBuffer
				}
			}
		}

		s.Editor = bufToOpen.Editors()[0]
		s.Minibuffer = E.Minibuffer.Editor()
		s.Focus(s.Editor)
		E.focus = bufToOpen.Editors()[0]

		if erroedOpening {
			s.Warning = "Failed to open files from command line."
		}
	})
}
