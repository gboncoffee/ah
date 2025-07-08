package main

import "errors"

type Minibuffer struct {
	prompt  string
	content string
	editor  *Editor
	commit  func(string)
	cancel  func()
}

func (b *Minibuffer) Init() {
	b.prompt = ""
	b.editor = NewEditor(b)
	b.editor.AddCursor(Cursor{Begin: 0, End: 1})
	b.editor.DefaultStyle = &E.Colors.Minibuffer
	b.editor.CursorStyle = &E.Colors.Cursor
}

func (b *Minibuffer) Editor() *Editor {
	return b.editor
}

func (b *Minibuffer) Insert(disp int, c byte) error {
	if c != '\n' {
		b.content = b.content[:disp] + string(c) + b.content[disp:]
		return nil
	}

	return errors.New("i'm the minibuffer")
}

func (b *Minibuffer) Get(disp int) (byte, error) {
	if disp < len(b.prompt) {
		return b.prompt[disp], nil
	}

	disp -= len(b.prompt)

	if len(b.content) >= disp {
		return 0, errors.New("nothing to read")
	}

	return b.content[disp], nil
}

func (b *Minibuffer) Delete(disp int) error {
	return nil
}
