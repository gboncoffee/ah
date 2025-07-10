package main

import (
	"math"
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
	e.vs = make(ui.VirtualEditorScreen, height)
	for i := range e.vs {
		e.vs[i] = make([]ui.VirtualRune, width)
	}
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

func (e *Editor) height() int {
	return len(e.vs)
}

func (e *Editor) width() int {
	return len(e.vs[0])
}

func (e *Editor) actualLineNumWidth() (w int) {
	if e.NumberColumn {
		w = max(int(math.Log10(float64(e.firstLine+e.height())))+1, 5)
	}
	return
}

func (e *Editor) fillBackground() {
	for i := range e.vs {
		for j := range e.vs[i] {
			e.vs[i][j].Style = e.DefaultStyle
			e.vs[i][j].Rune = ' '
		}
	}
}

func (e *Editor) fillLineNumStyle(width int) {
	if e.NumberColumn {
		for h := range e.height() {
			for i := range width {
				e.vs[h][i].Style = e.NumberColumnStyle
			}
		}
	}
}

func (e *Editor) maxWidth(lineNumWidth int) (w int) {
	w = e.width()
	if e.TextWidth != 0 {
		w = min(w, lineNumWidth+e.TextWidth+1)
	}
	return
}

func (e *Editor) drawLineNumber(line int, lineNumWidth int, cury int) {
	num := strconv.Itoa(line)
	curx := lineNumWidth - 1 - len(num)
	for _, c := range num {
		e.vs[cury][curx].Rune = c
		curx++
	}
}

func (e *Editor) drawCursors(focus bool, curx, cury, disp int) {
	if focus {
		for _, c := range e.cursors {
			if disp >= c.Begin && disp < c.End {
				e.vs[cury][curx].Style = e.CursorStyle
				break
			}
		}
	}
}

func (e *Editor) render(focus bool) {
	width := e.width()
	height := e.height()

	e.fillBackground()

	// Compute number column width.
	lineNumWidth := e.actualLineNumWidth()

	e.fillLineNumStyle(lineNumWidth)

	maxWidth := e.maxWidth(lineNumWidth)

	curx := 0
	cury := 0
	disp := e.disp
	line := e.firstLine
	continuing := false

	r, newDisp := getRune(e.buffer, disp)
	if disp != newDisp {
	line:
		for cury < height {
			if e.NumberColumn {
				if !continuing {
					line++
					e.drawLineNumber(line, lineNumWidth, cury)
				}
				curx = lineNumWidth
				continuing = false
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

					e.drawCursors(focus, curx, cury, disp)
					if r != '\t' {
						e.vs[cury][curx].Rune = r
						curx++
					} else {
						curx += 8
					}
				} else if r == '\n' {
					e.drawCursors(focus, curx, cury, disp)
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
