package main

import (
	"fmt"
	"math"
	"strconv"

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

	e.DefaultStyle = &E.Colors.Default
	e.NumberColumnStyle = &E.Colors.NumberColumn
	e.CursorStyle = &E.Colors.Cursor
	e.TextWidthColumnStyle = &E.Colors.TextWidthColumn

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

func (e *Editor) drawLineNumber(line int, lineNumWidth int, cury int) {
	num := strconv.Itoa(line + 1)
	curx := lineNumWidth - 1 - len(num)
	if curx < 0 {
		panic(fmt.Sprintf("lineNumWidth %v len(num) %v", lineNumWidth, len(num)))
	}
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
	height := e.height()
	width := e.width()

	e.fillBackground()

	// Compute number column width.
	lineNumWidth := e.actualLineNumWidth()
	e.fillLineNumStyle(lineNumWidth)

	curx := 0
	cury := 0
	disp := e.disp
	line := e.firstLine
	continuing := false

line:
	for cury < height {
		c, err := e.buffer.Get(disp)
		if err != nil {
			break line
		}

		// Draw line numbers.
		if e.NumberColumn && !continuing {
			e.drawLineNumber(line, lineNumWidth, cury)
		}
		curx += lineNumWidth
		continuing = false

		for {
			e.vs[cury][curx].Rune = c
			if e.TextWidth != 0 && curx-lineNumWidth == e.TextWidth {
				e.vs[cury][curx].Style = e.TextWidthColumnStyle
			}
			e.drawCursors(focus, curx, cury, disp)

			disp++

			switch c {
			case '\n':
				// Maybe draw the text width line at the end.
				if curx-lineNumWidth < e.TextWidth &&
					e.TextWidth != 0 && lineNumWidth+e.TextWidth < width {
					e.vs[cury][lineNumWidth+e.TextWidth].Style =
						e.TextWidthColumnStyle
				}
				cury++
				curx = 0
				line++
				continue line
			case '\t':
				curx = lineNumWidth + (((curx-lineNumWidth)+8)/8)*8
			default:
				curx++
			}

			if curx >= width {
				curx = 0
				cury++
				continuing = true
				continue line
			}

			c, err = e.buffer.Get(disp)
			if err != nil {
				cury++
				break line
			}
		}
	}

	// Fill eob.
	for cury < height {
		e.vs[cury][0].Rune = '~'
		e.vs[cury][0].Style = e.NumberColumnStyle
		cury++
	}
}
