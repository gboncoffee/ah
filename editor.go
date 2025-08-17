package main

import (
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
	cursors  []Cursor
	lineView LineView

	tabSize      int
	textWidth    int
	numberColumn bool

	defaultStyle         *tcell.Style
	numberColumnStyle    *tcell.Style
	cursorStyle          *tcell.Style
	textWidthColumnStyle *tcell.Style
}

func NewEditor(b Buffer) (e *Editor) {
	e = new(Editor)
	e.lineView = NewLineView(b)
	e.tabSize = 8

	e.defaultStyle = &E.Colors.Default
	e.numberColumnStyle = &E.Colors.NumberColumn
	e.cursorStyle = &E.Colors.Cursor
	e.textWidthColumnStyle = &E.Colors.TextWidthColumn

	return
}

func (e *Editor) actualLineNumWidth(height int) (w int) {
	if e.numberColumn {
		w = max(int(math.Log10(float64(e.lineView.FirstLine()+height)))+1, 5)
	}
	return
}

func (e *Editor) fill(
	width, height int,
	setRune ui.SetRuneFunc,
	setStyle ui.SetStyleFunc,
) {
	for y := range height {
		for x := range width {
			setRune(x, y, ' ')
			setStyle(x, y, e.defaultStyle)
		}
	}
}

func (e *Editor) renderLineNum(
	y int,
	l Line,
	n int,
	lnw int,
	setRune ui.SetRuneFunc,
	setStyle ui.SetStyleFunc,
) int {
	if l.continuation {
		if !e.numberColumn {
			return n
		}

		for x := range lnw - 1 {
			setStyle(x, y, e.numberColumnStyle)
			setRune(x, y, '>')
		}
		setStyle(lnw-1, y, e.numberColumnStyle)

		return n
	}

	if !e.numberColumn {
		return n + 1
	}

	num := strconv.Itoa(n + 1)
	x := 0
	for range lnw - 1 - len(num) {
		setStyle(x, y, e.numberColumnStyle)
		x++
	}
	for _, c := range num {
		setStyle(x, y, e.numberColumnStyle)
		setRune(x, y, c)
		x++
	}

	setStyle(x, y, e.numberColumnStyle)

	return n + 1
}

func (e *Editor) renderCursors(disp, x, y int, setStyle ui.SetStyleFunc) {
	for _, c := range e.cursors {
		if c.Begin <= disp && c.End > disp {
			setStyle(x, y, e.cursorStyle)
			return
		}
	}
}

func (e *Editor) renderTW(x, y, lnw int, setStyle ui.SetStyleFunc) {
	if e.textWidth == 0 {
		return
	}

	if x-lnw == e.textWidth {
		setStyle(x, y, e.textWidthColumnStyle)
	}
}

func (e *Editor) Render(
	width, height int,
	focus bool,
	setRune ui.SetRuneFunc,
	setStyle ui.SetStyleFunc,
) {
	e.fill(width, height, setRune, setStyle)

	lnw := e.actualLineNumWidth(height)
	actualWidth := width - lnw
	e.lineView.Update(height, actualWidth)

	line := e.lineView.FirstLine()
	x := 0
	y := 0
	disp := e.lineView.Disp()

	for _, l := range e.lineView.Lines() {
		line = e.renderLineNum(y, l, line, lnw, setRune, setStyle)
		x += lnw
		for _, c := range l.content {
			setRune(x, y, c.c)
			setStyle(x, y, e.defaultStyle)
			e.renderTW(x, y, lnw, setStyle)
			e.renderCursors(disp, x, y, setStyle)
			x += c.size
			disp++
		}
		// Render text width if applicable.
		if e.textWidth > 0 &&
			!l.continuation &&
			x-lnw < e.textWidth &&
			e.textWidth < actualWidth {

			setStyle(lnw+e.textWidth, y, e.textWidthColumnStyle)
		}
		x = 0
		y++
	}
	for y < height {
		setRune(0, y, '~')
		setStyle(0, y, e.numberColumnStyle)
		if lnw > 1 {
			for i := range lnw - 1 {
				setStyle(i+1, y, e.numberColumnStyle)
				i++
			}
		}
		y++
	}
}
