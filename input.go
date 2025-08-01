package main

import "github.com/gboncoffee/ah/ui"

func (e *Haza) Event(event ui.Event) {
	switch ev := event.(type) {
	case *ui.KeyPress:
		e.KeyPress(ev)
	case *ui.RuneEntered:
		e.RuneEntered(ev)
	case *ui.Resize:
		e.Ui.Update(func(_ *ui.State) {})
	}
}

func (e *Haza) KeyPress(key *ui.KeyPress) {
	switch key.Key {
	case ui.KeyCtrlQ:
		e.Ui.Exit()
	case ui.KeyCtrlS:
		e.Save()
	case ui.KeyCtrlZ:
		e.focus.Undo()
	case ui.KeyCtrlY:
		e.focus.Redo()
	case ui.KeyBackspace, ui.KeyBackspace2:
		//e.focus.CursorLeft()
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
