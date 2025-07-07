package main

import (
	"errors"

	"github.com/deadpixi/rope"
)

type FileBuffer struct {
	Name    string
	Editors []*Editor
	content rope.Rope
	lines   int
}

func (b *FileBuffer) NewEditor() (e *Editor) {
	e = new(Editor)
	e.buffer = b
	e.TextWidth = 80
	e.NumberColumn = true

	e.DefaultStyle = &COLORS.Default
	e.NumberColumnStyle = &COLORS.NumberColumn
	e.CursorStyle = &COLORS.Cursor
	e.TextWidthColumnStyle = &COLORS.TextWidthColumn

	e.AddCursor(Cursor{Begin: 0, End: 1})

	return
}

func (b *FileBuffer) Insert(disp int, c byte) error {
	b.content = b.content.InsertString(disp, string(c))
	return nil
}

func (b *FileBuffer) Get(disp int) (byte, error) {
	var buffer [1]byte
	n, err := b.content.ReadAt(buffer[:], int64(disp))
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, errors.New("nothing to read")
	}

	return buffer[0], nil
}
