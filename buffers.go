package main

import (
	"errors"
	"strconv"

	"github.com/deadpixi/rope"
)

type buffer struct {
	history       []rope.Rope
	content       rope.Rope
	name          string
	numberOfLines int
}

func (e *editor) addUnnamedBuffer() *buffer {
	b, err := e.addBuffer("\n", "unnamed/")
	if err == nil {
		return b
	}

	i := 1
	for err != nil {
		b, err = e.addBuffer("\n", "unnamed"+strconv.Itoa(i)+"/")
		i++
	}

	return b
}

func (e *editor) addBuffer(content string, name string) (*buffer, error) {
	for _, buf := range e.buffers {
		if buf.name == name {
			return nil, errors.New("buffer exists")
		}
	}

	buf := newBuffer(content, name)
	e.buffers = append(e.buffers, buf)
	return buf, nil
}

func newBuffer(content string, name string) *buffer {
	return &buffer{
		content: rope.NewString(content),
		name:    name,
	}
}

func (e *editor) getBufferByName(name string) (*buffer, error) {
	for _, b := range e.buffers {
		if b.name == name {
			return b, nil
		}
	}

	return nil, errors.New("no such buffer")
}

func (b *buffer) commit() {
	b.history = append(b.history, b.content)
}

func (b *buffer) insert(position int, c rune) {
	b.content = b.content.InsertString(position, string(c))
	if c == '\n' {
		b.numberOfLines++
	}
}
