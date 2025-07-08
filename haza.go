package main

import (
	"errors"
	"fmt"
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
	Warning         tcell.Style
}

var E Haza

type FocusMode int

type Haza struct {
	Buffers    []*FileBuffer
	Minibuffer Minibuffer
	Logbuffer  []string
	Editors    []*Editor
	Colors     Colorscheme
	Ui         *ui.Ui
	focus      *Editor
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

	e.Colors.Warning = e.Colors.Minibuffer.Foreground(tcell.ColorRed)
	e.Colors.Message = e.Colors.Minibuffer
}

func (e *Haza) InitMinibuffer() {
	e.Minibuffer.Init()
}

func (e *Haza) Log(msg string) {
	e.Logbuffer = append(e.Logbuffer, msg)
}

func (e *Haza) LogError(msg string) {
	e.Logbuffer = append(e.Logbuffer, fmt.Sprintf("error: %v", msg))
}

func (e *Haza) Open(file string) (*FileBuffer, error) {
	f, err := os.Open(file)
	if err != nil {
		if os.IsPermission(err) {
			return nil, err
		}

		// Create an empty buffer as the file does not exists.
		return e.NewBuffer("", file), nil
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
