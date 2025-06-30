package main

import (
	"github.com/gboncoffee/ah/ui"

	"github.com/deadpixi/rope"
	"github.com/gdamore/tcell/v2"
)

type Colorscheme struct {
	Default         tcell.Style
	Cursor          tcell.Style
	NumberColumn    tcell.Style
	TextWidthColumn tcell.Style
	Minibuffer      tcell.Style
}

var COLORS Colorscheme

type Haza struct {
	ui         *ui.Ui
	Buffers    []*FileBuffer
	Minibuffer Minibuffer
	Editors    []*Editor
}

func (e *Haza) NewBuffer(content string, name string) *FileBuffer {
	buf := new(FileBuffer)
	*buf = FileBuffer{
		Name:    name,
		content: rope.NewString(content),
		Editors: []*Editor{buf.NewEditor()},
	}
	e.Buffers = append(e.Buffers, buf)

	return buf
}

func (e *Haza) InitColors() {
	COLORS.Default = tcell.StyleDefault.
		Foreground(tcell.ColorReset).Background(tcell.ColorReset)

	COLORS.Cursor = COLORS.Default.Reverse(true)
	COLORS.NumberColumn = COLORS.Default.
		Background(tcell.ColorGray)

	COLORS.TextWidthColumn = COLORS.Default.
		Foreground(tcell.ColorReset).Background(tcell.ColorSilver)

	COLORS.Minibuffer = COLORS.NumberColumn
}

func (e *Haza) InitMinibuffer() {
	e.Minibuffer.Init()
}
