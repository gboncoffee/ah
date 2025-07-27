package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type VirtualRune struct {
	Style *tcell.Style
	Rune  rune
}

type VirtualEditorScreen [][]VirtualRune

type EditorView interface {
	VirtualScreen(width, height int, focus bool) VirtualEditorScreen
}

func (ui *Ui) renderEditor(e EditorView, x, y, width, height int, focus bool) {
	vs := e.VirtualScreen(width, height, focus)
	for i := range vs {
		for j := range vs[i] {
			if vs[i][j].Style == nil {
				panic(fmt.Sprintf("style at %v %v", i, j))
			}
			ui.screen.SetContent(x+j, y+i, vs[i][j].Rune, nil, *vs[i][j].Style)
		}
	}
}
