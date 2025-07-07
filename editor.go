package main

import (
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
	firstLine            int
	disp                 int
	vs                   ui.VirtualEditorScreen
	TextWidth            int
	NumberColumn         bool
	DefaultStyle         *tcell.Style
	NumberColumnStyle    *tcell.Style
	CursorStyle          *tcell.Style
	TextWidthColumnStyle *tcell.Style
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
		}
	}
	e.cursors = slices.Insert(e.cursors, i, cursor)
}

func (e *Editor) RuneEntered(re *ui.RuneEntered) {
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

	UI.Update(func(_ *ui.State) {})
}

func (e *Editor) KeyEvent(key *ui.KeyPress) {
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
	end := false

	r, newDisp := getRune(e.buffer, disp)
	if newDisp == disp {
		end = true
	}

line:
	for cury < height {
		// Render line numbers or line continuations.
		if e.NumberColumn {
			if continuing {
				for curx < lineNumWidth {
					e.vs[cury][curx].Style = e.NumberColumnStyle
					curx++
				}
			} else if end {
				e.vs[cury][curx].Style = e.NumberColumnStyle
				e.vs[cury][curx].Rune = '~'
				curx++
				for curx < lineNumWidth {
					e.vs[cury][curx].Style = e.NumberColumnStyle
					curx++
				}
				curx = 0
				cury++
				continue line
			} else {
				num := strconv.Itoa(line)
				for curx < lineNumWidth-1-len(num) {
					e.vs[cury][curx].Style = e.NumberColumnStyle
					e.vs[cury][curx].Rune = ' '
					curx++
				}
				for _, c := range num {
					e.vs[cury][curx].Style = e.NumberColumnStyle
					e.vs[cury][curx].Rune = c
					curx++
				}
				e.vs[cury][curx].Style = e.NumberColumnStyle
				curx++
			}
		}

		for curx < maxWidth {
			// Render cursors.
			if focus {
				for _, cursor := range e.cursors {
					if disp >= cursor.Begin && disp < cursor.End {
						e.vs[cury][curx].Style = e.CursorStyle
					}
				}
			}

			if r == '\n' {
				line++
				cury++
				curx = 0
				continuing = false

				disp = newDisp
				r, newDisp = getRune(e.buffer, disp)
				if newDisp == disp {
					end = true
				}

				continue line
			}

			if r != 0 {
				e.vs[cury][curx].Rune = r
			}

			curx++

			disp = newDisp
			r, newDisp = getRune(e.buffer, disp)
			if newDisp == disp {
				end = true
			}
		}

		continuing = true
		curx = 0
		cury++
	}

	// Fill text width column.
	if maxWidth < width {
		for h := range height {
			e.vs[h][maxWidth].Style = e.TextWidthColumnStyle
		}
	}
}
