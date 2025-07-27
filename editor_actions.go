package main

import (
	"fmt"
	"slices"

	"github.com/gboncoffee/ah/ui"
)

// TODO: rewrite most of this with the new rune-addressable rope.

func (e *Editor) Save() {
	if fb, ok := e.buffer.(*FileBuffer); ok {
		E.Ui.Update(func(s *ui.State) {
			s.Message = fmt.Sprintf("Saving buffer %v...", fb.Name())
		})
		if err := fb.TrySave(); err != nil {
			E.Ui.Update(func(s *ui.State) {
				s.Warning = fmt.Sprintf("Error: %v", err)
			})
		} else {
			E.Ui.Update(func(s *ui.State) {
				s.Message = fmt.Sprintf("Saved buffer %v", fb.Name())
			})
		}
	} else {
		E.Ui.Update(func(s *ui.State) {
			s.Warning = "Error: cannot save non-file buffer."
		})
	}
}

//
// Most (actual) functions are inside update as to not cause race conditions.
//

func (e *Editor) AddCursor(cursor Cursor) {
	i := 0
	for i < len(e.cursors) {
		if e.cursors[i].Begin >= cursor.Begin {
			e.cursors = slices.Insert(e.cursors, i, cursor)
			return
		}
	}
	e.cursors = slices.Insert(e.cursors, i, cursor)
}

func (e *Editor) RuneEntered(re rune) {
	// The function is "utf8-unsafe" in the sense we're not strictly checking
	// if the cursors will be normalized and stuff like that but I'm pretty sure
	// this works. It's like using unsafe Rust. Or C++ without condoms.
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			e.buffer.Insert(e.cursors[i].Begin, re)
			for j := range e.cursors[i:] {
				e.cursors[j].Begin++
				e.cursors[j].End++
			}
		}
	})
}

func (e *Editor) Delete() {
	// This function is also "utf8-unsafe".
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			delRange := e.cursors[i].End - e.cursors[i].Begin
			for range delRange {
				e.buffer.Delete(e.cursors[i].Begin - 1)
			}
			for j := range e.cursors[i:] {
				e.cursors[j].Begin -= delRange
				e.cursors[j].End -= delRange
			}
		}
	})
}

func (e *Editor) CursorLeft() {
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			e.cursors[i] = e.cursorLeft(e.cursors[i])
		}
	})
}

func (e *Editor) cursorLeft(c Cursor) Cursor {
	if c.Begin > 0 {
		c.Begin--
	}
	c.End = c.Begin + 1

	return c
}

func (e *Editor) CursorRight() {
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			e.cursors[i] = e.cursorRight(e.cursors[i])
		}
	})
}

func (e *Editor) cursorRight(c Cursor) Cursor {
	if c.Begin+1 < e.buffer.Size() {
		c.Begin++
	}
	c.End = c.Begin + 1

	return c
}

func (e *Editor) gotoLineBegin(c Cursor) Cursor {
	r, err := e.buffer.Get(c.Begin - 1)
	for err == nil && r != '\n' {
		c.Begin--
		r, err = e.buffer.Get(c.Begin - 1)
	}
	c.End = c.Begin + 1
	return c
}

func (e *Editor) goRightUntilDispOrEol(c Cursor, disp int) Cursor {
	return c
}

func (e *Editor) CursorUp() {
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			e.cursors[i] = e.cursorUp(e.cursors[i])
		}
	})
}

func (e *Editor) cursorUp(c Cursor) Cursor {
	return c
}

func (e *Editor) CursorDown() {
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			e.cursors[i] = e.cursorDown(e.cursors[i])
		}
	})
}

func (e *Editor) cursorDown(c Cursor) Cursor {
	return c
}
