package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/gboncoffee/ah/ui"
)

func (e *Editor) RuneEntered(re rune) {
	// The function is inside the Update as to not cause race conditions in the
	// insertion.
	E.Ui.Update(func(_ *ui.State) {
		var buffer [4]byte
		slice := buffer[:]
		size := utf8.EncodeRune(slice, re)

		for _, byte := range slice[:size] {
			for i := range e.cursors {
				e.buffer.Insert(e.cursors[i].Begin, byte)
				for j := range e.cursors[i:] {
					e.cursors[j].Begin++
					e.cursors[j].End++
				}
			}
		}
	})
}

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

func (e *Editor) Backspace() {
	// The function is inside the Update as to not cause race conditions in the
	// deletion.
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			if e.cursors[i].Begin <= 0 {
				continue
			}
			err := e.buffer.Delete(e.cursors[i].Begin - 1)
			if err != nil {
				continue
			}
			for j := range e.cursors[i:] {
				e.cursors[j].Begin--
				e.cursors[j].End--
			}
		}
	})
}

func (e *Editor) CursorLeft() {
	// The function is inside the Update as to not cause race conditions in the
	// movement (yes that can actually happen).
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			if e.cursors[i].Begin >= 0 {
				// The cursor collapses on movement.
				e.cursors[i].Begin--
				e.cursors[i].End = e.cursors[i].Begin + 1
			}
		}
	})
}

func (e *Editor) CursorRight() {
	// The function is inside the Update as to not cause race conditions in the
	// movement (yes that can actually happen).
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			// Check if there's something remaining in the buffer.
			if _, err := e.buffer.Get(e.cursors[i].Begin + 1); err != nil {
				// The cursor collapses on movement.
				e.cursors[i].Begin++
				e.cursors[i].End = e.cursors[i].Begin + 1
			}
		}
	})
}
