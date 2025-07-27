package main

import (
	"errors"
	"slices"
)

type Minibuffer struct {
	prompt  []rune
	content []rune
	editor  *Editor
	commit  func(string)
	cancel  func()
}

func (b *Minibuffer) Init() {
	b.editor = NewEditor(b)
	b.editor.AddCursor(Cursor{Begin: 0, End: 1})
	b.editor.DefaultStyle = &E.Colors.Minibuffer
	b.editor.CursorStyle = &E.Colors.Cursor

	b.Reset()
}

func (b *Minibuffer) Reset() {
	b.prompt = []rune{}
	b.content = []rune{'\n'}
}

func (b *Minibuffer) Editor() *Editor {
	return b.editor
}

func (b *Minibuffer) Insert(idx int, r rune) error {
	if r != '\n' {
		b.content = slices.Insert(b.content, idx+len(b.prompt), r)
		return nil
	}

	return errors.New("i'm the minibuffer")
}

func (b *Minibuffer) Get(idx int) (rune, error) {
	if idx < len(b.prompt) {
		return b.prompt[idx], nil
	}

	idx -= len(b.prompt)

	if len(b.content) <= idx {
		return 0, errors.New("nothing to read")
	}

	return b.content[idx], nil
}

func (b *Minibuffer) Delete(disp int) error {
	return nil
}

func (b *Minibuffer) Size() int {
	return len(b.prompt) + len(b.content)
}
