package main

import (
	"errors"
)

type Minibuffer struct {
	Prompt  string
	Editor  *Editor
	Content string
}

func (m *Minibuffer) Init() {
	m.Prompt = ""
	m.Editor = &Editor{
		buffer:       m,
		cursors:      []Cursor{{Begin: 0, End: 1}},
		NumberColumn: false,
		DefaultStyle: &COLORS.Minibuffer,
		CursorStyle:  &COLORS.Cursor,
	}
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
