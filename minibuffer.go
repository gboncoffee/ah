package main

import (
	"errors"

	pt "github.com/gboncoffee/gopiecetable"
)

var ErrorImTheMinibuffer = errors.New("i'm the minibuffer")
var ErrorInPrompt = errors.New("in the prompt")

type Minibuffer struct {
	buffer     *pt.PieceTable[rune]
	promptSize int
	editor     *Editor
	commit     func(string)
	cancel     func()
}

func (b *Minibuffer) Init() {
	b.Reset()
	b.editor.AddCursor(Cursor{Begin: 0, End: 1})
	b.editor.DefaultStyle = &E.Colors.Minibuffer
	b.editor.CursorStyle = &E.Colors.Cursor

	b.Reset()
}

func (b *Minibuffer) Reset() {
	b.buffer = pt.FromString("\n")
	b.editor = NewEditor(b)
}

func (b *Minibuffer) Editor() *Editor {
	return b.editor
}

func (b *Minibuffer) Insert(idx int, r rune) error {
	if r != '\n' {
		return b.buffer.Insert(idx, r)
	}
	if idx < b.promptSize {
		return ErrorInPrompt
	}

	return ErrorImTheMinibuffer
}

func (b *Minibuffer) Get(idx int) (rune, error) {
	return b.buffer.Get(idx)
}

func (b *Minibuffer) Delete(idx int) error {
	if idx < b.promptSize {
		return ErrorInPrompt
	}

	return b.buffer.Delete(idx)
}

func (b *Minibuffer) Size() int {
	return b.buffer.Size()
}

func (b *Minibuffer) Undo() (int, error) {
	return b.buffer.Undo()
}

func (b *Minibuffer) Redo() (int, error) {
	return b.buffer.Redo()
}

func (b *Minibuffer) Commit() {
	b.commit(pt.String(b.buffer)[b.promptSize:])
}

func (b *Minibuffer) Cancel() {
	b.cancel()
}
