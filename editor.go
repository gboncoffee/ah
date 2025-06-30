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
	buffer               Buffer
	cursors              []Cursor
	firstLine            int
	vs                   ui.VirtualEditorScreen
	width                int
	height               int
	TextWidth            int
	numColumnSize        int
	NumberColumn         bool
	DefaultStyle         *tcell.Style
	NumberColumnStyle    *tcell.Style
	CursorStyle          *tcell.Style
	TextWidthColumnStyle *tcell.Style
}

func (e *Editor) Event(ev ui.Event) {}

func (e *Editor) VirtualScreen(width, height int, focus bool) ui.VirtualEditorScreen {
	if e.vs != nil && width == e.width && height == e.height {
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
	e.width = width
	e.height = height
}

// TODO: Make the render function less ugly. Holy fucking shit. It's just too
// ugly. Please.

func (e *Editor) render(focus bool) {
	disp := e.buffer.Displacement(e.firstLine)
	reader := e.buffer.DisplacedReader(disp)

	for i := range e.height {
		for j := range e.width {
			e.vs[i][j] = ui.VirtualRune{
				Rune:  ' ',
				Style: e.DefaultStyle,
			}
		}
	}

	if e.NumberColumn {
		e.numColumnSize =
			max(int(math.Floor(math.Log10(float64(e.buffer.Lines()))))+2, 5)
	}

	// Fill textwidth column if needed.
	maxWidth := e.width
	if e.TextWidth != 0 && e.TextWidth+e.numColumnSize < e.width {
		maxWidth = e.TextWidth + e.numColumnSize
		for i := range e.height {
			e.vs[i][maxWidth].Style = e.TextWidthColumnStyle
		}
	}

	curx := 0
	cury := 0
	curline := e.firstLine + 1
	stopReading := false

	for cury < e.height {
		num := "~ "
		var line string
		if !stopReading {
			var err error
			line, err = reader.ReadString('\n')
			if err != nil {
				stopReading = true
			} else if e.NumberColumn {
				num = strconv.Itoa(curline) + " "
			}
		}

		if e.NumberColumn {
			for range e.numColumnSize - 2 {
				num = " " + num
			}

			for curx+len(num) < e.numColumnSize {
				e.vs[cury][curx].Style = e.NumberColumnStyle
				curx++
			}
			for _, c := range num {
				e.vs[cury][curx].Style = e.NumberColumnStyle
				e.vs[cury][curx].Rune = c
				curx++
			}
		}

		for _, c := range line {
			color := e.DefaultStyle
			if focus {
				for _, cursor := range e.cursors {
					if cursor.Begin <= disp && cursor.End > disp {
						color = e.CursorStyle
					}
				}
			}
			e.vs[cury][curx].Style = color
			e.vs[cury][curx].Rune = c
			disp++
			curx++
			if curx >= maxWidth {
				cury++
				// Dijkstra probably hates me.
				if cury > e.height {
					return
				}
				curx = 0
				if e.NumberColumn {
					for curx < e.numColumnSize-1 {
						e.vs[cury][curx].Style = e.NumberColumnStyle
						e.vs[cury][curx].Rune = '>'
						curx++
					}
					e.vs[cury][curx].Style = e.NumberColumnStyle
					curx++
				}
			}
		}

		curline++
		cury++
		curx = 0
	}
}
