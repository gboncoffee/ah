package main

import (
	"bufio"
	"strings"
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

func (m *Minibuffer) Displacement(disp int) int {
	return 0
}

func (m *Minibuffer) DisplacedReader(disp int) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(m.Content))
}

func (m *Minibuffer) Insert(disp int, content string) {
	m.Content = m.Content[:disp] + content + m.Content[disp:]
}

func (m *Minibuffer) Lines() int {
	return 1
}
