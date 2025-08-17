package ui

import (
	"github.com/gdamore/tcell/v2"
)

type SetRuneFunc func(x, y int, c rune)
type SetStyleFunc func(x, y int, s *tcell.Style)

type EditorView interface {
	Render(
		width, height int,
		focus bool,
		setRune SetRuneFunc,
		setStyle SetStyleFunc,
	)
}

type cell struct {
	c rune
	s *tcell.Style
}

func (ui *Ui) renderEditor(
	e EditorView,
	xb, yb, width, height int,
	focus bool,
) {
	vs := make([][]cell, height)
	for i := range vs {
		vs[i] = make([]cell, width)
	}

	e.Render(width, height, focus, func(x, y int, c rune) {
		vs[y][x].c = c
	}, func(x, y int, s *tcell.Style) {
		vs[y][x].s = s
	})

	for y, l := range vs {
		for x, c := range l {
			ui.screen.SetContent(xb+x, yb+y, c.c, nil, *c.s)
		}
	}
}
