package main

import "errors"

type Minibuffer struct {
	Prompt  string
	Editor  *Editor
	Content string
}

func (m *Minibuffer) Init() {
	m.Prompt = ""
	m.Editor = NewEditor(m)
	m.Editor.AddCursor(Cursor{Begin: 0, End: 1})
	m.Editor.DefaultStyle = &E.Colors.Minibuffer
	m.Editor.CursorStyle = &E.Colors.Cursor
}

func (b *Minibuffer) Insert(disp int, c byte) error {
	if c != '\n' {
		b.Content = b.Content[:disp] + string(c) + b.Content[disp:]
	}

	return nil
}

func (b *Minibuffer) Get(disp int) (byte, error) {
	if len(b.Content) >= disp {
		return 0, errors.New("nothing to read")
	}

	return b.Content[disp], nil
}

func (b *Minibuffer) Delete(disp int) error {
	if disp >= len(b.Content) {
		return errors.New("index out of range")
	}

	if disp < len(b.Content)-1 {
		b.Content = b.Content[:disp] + b.Content[disp:]
		return nil
	}

	b.Content = b.Content[:disp]

	return nil
}
