package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/deadpixi/rope"
)

type FileBuffer struct {
	name    string
	editors []*Editor
	content rope.Rope
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
	fb.content = rope.NewString(content)

	for _, c := range content {
		if c == '\n' {
			fb.lines++
		}
	}

	return
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

	_, err = f.WriteString(b.content.String())
	if err != nil {
		return fmt.Errorf("cannot write to file %v: %v", b.name, err)
	}

	return nil
}

func (b *FileBuffer) NewEditor() (e *Editor) {
	e = NewEditor(b)

	e.TextWidth = 80
	e.NumberColumn = true

	e.DefaultStyle = &E.Colors.Default
	e.NumberColumnStyle = &E.Colors.NumberColumn
	e.CursorStyle = &E.Colors.Cursor
	e.TextWidthColumnStyle = &E.Colors.TextWidthColumn

	e.AddCursor(Cursor{Begin: 0, End: 1})

	b.editors = append(b.editors, e)

	return
}

func (b *FileBuffer) Insert(disp int, c byte) error {
	if c == '\n' {
		b.lines++
	}
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

func (b *FileBuffer) Delete(disp int) error {
	b.content = b.content.Delete(disp, 1)

	return nil
}
