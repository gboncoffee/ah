package main

import (
	"bufio"
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

	return
}

func (b *FileBuffer) Displacement(line int) (disp int) {
	reader := bufio.NewReader(b.content.Reader())
	for nlines := 0; nlines != line; nlines++ {
		s, _ := reader.ReadString('\n')
		disp += len(s)
	}
	return
}

func (b *FileBuffer) DisplacedReader(disp int) (reader *bufio.Reader) {
	return bufio.NewReader(b.content.OffsetReader(disp))
}

func (b *FileBuffer) Lines() int {
	return b.lines
}

func (b *FileBuffer) Insert(disp int, content string) {
	b.content = b.content.InsertString(disp, content)
}
