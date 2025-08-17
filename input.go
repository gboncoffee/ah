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
	case *ui.MouseScrollDown:
		e.MouseScrollDown()
	case *ui.MouseScrollUp:
		e.MouseScrollUp()
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
	case ui.KeyDown:
		e.focus.CursorDown()
	case ui.KeyUp:
		e.focus.CursorUp()
	case ui.KeyPgDn:
		e.PageDown()
	case ui.KeyPgUp:
		e.PageUp()
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

func (e *Haza) MouseScrollDown() {
	if e.focus != e.Minibuffer.Editor() {
		e.Ui.Update(func(_ *ui.State) {
			e.focus.ScrollForward(1)
		})
	}
}

func (e *Haza) MouseScrollUp() {
	if e.focus != e.Minibuffer.Editor() {
		e.Ui.Update(func(_ *ui.State) {
			e.focus.ScrollBack(1)
		})
	}
}

func (e *Haza) PageDown() {
	if e.focus != e.Minibuffer.Editor() {
		e.Ui.Update(func(state *ui.State) {
			e.focus.ScrollForward(state.EditorHeight() / 2)
		})
	}
}

func (e *Haza) PageUp() {
	if e.focus != e.Minibuffer.Editor() {
		e.Ui.Update(func(state *ui.State) {
			e.focus.ScrollBack(state.EditorHeight() / 2)
		})
	}
}
