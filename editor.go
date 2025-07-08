package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/gboncoffee/ah/ui"

	"github.com/gdamore/tcell/v2"
)

type Cursor struct {
	Begin int
	End   int
}

// Implements the EditorView ui interface.
type Editor struct {
	buffer               Buffer
	cursors              []Cursor
	masterCursor         *Cursor
	firstLine            int
	disp                 int
	vs                   ui.VirtualEditorScreen
	TabSize              int
	TextWidth            int
	NumberColumn         bool
	DefaultStyle         *tcell.Style
	NumberColumnStyle    *tcell.Style
	CursorStyle          *tcell.Style
	TextWidthColumnStyle *tcell.Style
}

func NewEditor(b Buffer) (e *Editor) {
	e = new(Editor)
	e.buffer = b
	e.TabSize = 8

	return
}

func (e *Editor) Event(event ui.Event) {
	switch ev := event.(type) {
	case *ui.KeyPress:
		e.KeyEvent(ev)
	case *ui.RuneEntered:
		e.RuneEntered(ev)
	}
}

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

func (e *Editor) RuneEntered(re *ui.RuneEntered) {
	// The function is inside the Update as to not cause race conditions in the
	// insertion.
	E.Ui.Update(func(_ *ui.State) {
		var buffer [4]byte
		slice := buffer[:]
		size := utf8.EncodeRune(slice, re.Rune)

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

func (e *Editor) KeyEvent(key *ui.KeyPress) {
	switch key.Key {
	case ui.KeyCtrlS:
		e.save()
	case ui.KeyBackspace, ui.KeyBackspace2:
		e.backspace()
	}
}

func (e *Editor) save() {
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

func (e *Editor) backspace() {
	// The function is inside the Update as to not cause race conditions in the
	// deletion.
	E.Ui.Update(func(_ *ui.State) {
		for i := range e.cursors {
			if e.cursors[i].Begin > 0 {
				e.buffer.Delete(e.cursors[i].Begin - 1)
				for j := range e.cursors[i:] {
					e.cursors[j].Begin--
					e.cursors[j].End--
				}
			}
		}
	})
}

func (e *Editor) VirtualScreen(width, height int, focus bool) ui.VirtualEditorScreen {
	if e.vs != nil && len(e.vs[0]) == width && len(e.vs) == height {
		e.render(focus)
		return e.vs
	}
	e.newVS(width, height)
	e.render(focus)
	return e.vs
}

func (e *Editor) newVS(width, height int) {
	vs := make(ui.VirtualEditorScreen, height)
	for i := range vs {
		vs[i] = make([]ui.VirtualRune, width)
	}

	e.vs = vs
}

// TODO: Make the render function less ugly. Holy fucking shit. It's just too
// ugly. Please.

func getRune(b Buffer, disp int) (r rune, newDisp int) {
	var buffer [4]byte
	firstByte, err := b.Get(disp)
	if err != nil {
		return ' ', disp
	}
	buffer[0] = firstByte
	buffer[1], _ = b.Get(disp + 1)
	buffer[2], _ = b.Get(disp + 2)
	buffer[3], _ = b.Get(disp + 3)

	var size int
	r, size = utf8.DecodeRune(buffer[:])
	return r, disp + size
}

func (e *Editor) render(focus bool) {
	width := len(e.vs[0])
	height := len(e.vs)

	// Fill default style and character.
	for i := range e.vs {
		for j := range e.vs[i] {
			e.vs[i][j].Style = e.DefaultStyle
			e.vs[i][j].Rune = ' '
		}
	}

	// Compute number column width.
	lineNumWidth := 0
	if e.NumberColumn {
		lineNumWidth = max(int(math.Log10(float64(e.firstLine+height)))+1, 5)
	}

	// Set style for number column.
	if e.NumberColumn {
		for h := range height {
			for i := range lineNumWidth {
				e.vs[h][i].Style = e.NumberColumnStyle
			}
		}
	}

	// Compute max width
	maxWidth := width
	if e.TextWidth != 0 {
		maxWidth = min(width, lineNumWidth+e.TextWidth+1)
	}

	curx := 0
	cury := 0
	disp := e.disp
	line := e.firstLine
	continuing := false

	r, newDisp := getRune(e.buffer, disp)
	if disp != newDisp {
	line:
		for cury < height {
			// Draw line numbers.
			if e.NumberColumn {
				if !continuing {
					line++
					num := strconv.Itoa(line)
					curx += lineNumWidth - 1 - len(num)
					for _, c := range num {
						e.vs[cury][curx].Rune = c
						curx++
					}
					curx++
				} else {
					continuing = false
					curx += lineNumWidth
				}
			}

			for {
				if r != '\n' && r != 0 {
					if curx >= maxWidth {
						curx = 0
						cury++
						// Next iteration we'll use the same rune.
						continuing = true
						continue line
					}

					// Draw cursors.
					if focus {
						for _, c := range e.cursors {
							if disp >= c.Begin && disp < c.End {
								e.vs[cury][curx].Style = e.CursorStyle
								break
							}
						}
					}
					if r != '\t' {
						e.vs[cury][curx].Rune = r
						curx++
					} else {
						curx += 8
					}
				} else if r == '\t' {
					curx += 8
				} else if r == '\n' {
					// Draw cursors.
					if curx < width && focus {
						for _, c := range e.cursors {
							if disp >= c.Begin && disp < c.End {
								e.vs[cury][curx].Style = e.CursorStyle
								break
							}
						}
					}
					curx = 0
					cury++
					disp = newDisp
					r, newDisp = getRune(e.buffer, disp)
					if disp == newDisp {
						break line
					}
					continue line
				} else {
					curx++
				}

				disp = newDisp
				r, newDisp = getRune(e.buffer, disp)
				if disp == newDisp {
					break line
				}
			}
		}
	}

	// Fill remaining lines with line ending.
	for cury < height {
		e.vs[cury][0].Rune = '~'
		cury++
	}

	// Fill text width column.
	if maxWidth < width {
		for h := range height {
			e.vs[h][maxWidth].Style = e.TextWidthColumnStyle
		}
	}
}
