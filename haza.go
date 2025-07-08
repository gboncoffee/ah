package main

import (
	"errors"
	"io"
	"os"
	"unicode/utf8"

	"github.com/gboncoffee/ah/ui"

	"github.com/gdamore/tcell/v2"
)

type Colorscheme struct {
	Default         tcell.Style
	Cursor          tcell.Style
	NumberColumn    tcell.Style
	TextWidthColumn tcell.Style
	Minibuffer      tcell.Style
	Message         tcell.Style
}

var E Haza

type Haza struct {
	Buffers    []*FileBuffer
	Minibuffer Minibuffer
	Logbuffer  []string
	Editors    []*Editor
	Colors     Colorscheme
	Ui         *ui.Ui
}

func (e *Haza) NewBuffer(content string, name string) *FileBuffer {
	buf := NewFileBuffer(name, content)
	buf.NewEditor()

	e.Buffers = append(e.Buffers, buf)

	return buf
}

func (e *Haza) InitColors() {
	e.Colors.Default = tcell.StyleDefault.
		Foreground(tcell.ColorReset).Background(tcell.ColorReset)

	e.Colors.Cursor = e.Colors.Default.Reverse(true)
	e.Colors.NumberColumn = e.Colors.Default.
		Background(tcell.ColorGray)

	e.Colors.TextWidthColumn = e.Colors.Default.
		Foreground(tcell.ColorReset).Background(tcell.ColorSilver)

	e.Colors.Minibuffer = e.Colors.NumberColumn

	e.Colors.Message = e.Colors.Minibuffer.Foreground(tcell.ColorRed)
}

func (e *Haza) InitMinibuffer() {
	e.Minibuffer.Init()
}

func (e *Haza) Log(msg string) {
	e.Logbuffer = append(e.Logbuffer, msg)
}

func (e *Haza) Open(file string) (*FileBuffer, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if !utf8.Valid(content) {
		return nil, errors.New("file is not UTF-8")
	}

	return e.NewBuffer(string(content), file), nil
}
