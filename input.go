package main

import "github.com/gboncoffee/ah/ui"

func Event(e *Haza, event ui.Event) {
	switch ev := event.(type) {
	case *ui.KeyPress:
		e.KeyPress(ev)
	case *ui.RuneEntered:
		e.RuneEntered(ev)
	}
}

func (e *Haza) KeyPress(key *ui.KeyPress) {
	switch key.Key {
	case ui.KeyCtrlQ:
		e.Ui.Exit()
	case ui.KeyCtrlS:
		e.Save()
	case ui.KeyBackspace, ui.KeyBackspace2:
		e.focus.CursorLeft()
		e.focus.Delete()
	case ui.KeyLeft:
		e.focus.CursorLeft()
	case ui.KeyRight:
		e.focus.CursorRight()
	}
}

func (e *Haza) RuneEntered(key *ui.RuneEntered) {
	if key.Rune != '\n' || e.focus != e.Minibuffer.Editor() {
		e.focus.RuneEntered(key.Rune)
	}
}

func (e *Haza) Save() {
	if e.focus != e.Minibuffer.Editor() {
		e.focus.Save()
	}
}
