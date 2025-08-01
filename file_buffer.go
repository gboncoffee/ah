package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gboncoffee/ah/buffer"
)

type FileBuffer struct {
	name    string
	editors []*Editor
	content *buffer.Buffer[rune]
	lines   int
	ioLock  sync.Mutex
}

func (fb *FileBuffer) Name() string {
	return fb.name
}

func NewFileBuffer(name, content string) (fb *FileBuffer) {
	if len(content) < 1 || content[len(content)-1] != '\n' {
		content += string('\n')
	}

	fb = new(FileBuffer)
	fb.name = name
	fb.content = buffer.FromString(content)

	for _, c := range content {
		if c == '\n' {
			fb.lines++
		}
	}

	return
}

func (b *FileBuffer) Editors() []*Editor {
	return b.editors
}

func (b *FileBuffer) TrySave() error {
	if strings.HasPrefix(b.name, "//") {
		return fmt.Errorf("cannot save virtual file %v", b.name)
	}

	b.ioLock.Lock()
	defer b.ioLock.Unlock()

	f, err := os.Create(b.name)
	if err != nil {
		return fmt.Errorf("cannot save file %v: %v", b.name, err)
	}

	_, err = f.WriteString(buffer.String(b.content))
	if err != nil {
		return fmt.Errorf("cannot write to file %v: %v", b.name, err)
	}

	return nil
}

func (b *FileBuffer) NewEditor() (e *Editor) {
	e = NewEditor(b)

	e.TextWidth = 80
	e.NumberColumn = true

	e.AddCursor(Cursor{Begin: 0, End: 1})

	b.editors = append(b.editors, e)

	return
}

func (b *FileBuffer) Insert(idx int, r rune) error {
	if r == '\n' {
		b.lines++
	}
	b.content.Insert(idx, r)
	return nil
}

func (b *FileBuffer) Get(idx int) (rune, error) {
	return b.content.Get(idx)
}

func (b *FileBuffer) Delete(idx int) error {
	return b.content.Delete(idx)
}

func (b *FileBuffer) Size() int {
	return b.content.Size()
}

func (b *FileBuffer) Undo() (int, error) {
	return b.content.Undo()
}

func (b *FileBuffer) Redo() (int, error) {
	return b.content.Redo()
}
