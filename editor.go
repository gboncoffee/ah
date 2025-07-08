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
