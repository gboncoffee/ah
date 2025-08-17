package main

import "errors"

type Char struct {
	c    rune
	size int
}

type Line struct {
	content      []Char
	continuation bool
}

func (l Line) Size() (size int) {
	for _, c := range l.content {
		size += c.size
	}
	return
}

type LineView struct {
	buffer        Buffer
	firstLine     int
	disp          int
	sum           int
	width, height int
	lines         []Line
}

func NewLineView(b Buffer) LineView {
	return LineView{
		buffer:    b,
		firstLine: 0,
	}
}

func (l *LineView) Buffer() Buffer {
	return l.buffer
}

func (l *LineView) FirstLine() int {
	return l.firstLine
}

func (l *LineView) Disp() int {
	return l.disp
}

func (l *LineView) Size() (s int) {
	for _, l := range l.lines {
		s += len(l.content)
	}
	return
}

func (l *LineView) Update(height, width int) {
	l.lines = make([]Line, 0, height)
	idx := l.disp
	cline := -1
	cont := false
	l.sum = 0
	l.width = width
	l.height = height

	for cline < height-1 {
		// Hack to quit on the end.
		_, err := l.buffer.Get(idx)
		if err != nil {
			return
		}

		l.lines = append(l.lines, Line{})
		lineLength := 0
		cline++
		l.lines[cline].continuation = cont
		for lineLength < width {
			c, err := l.buffer.Get(idx)
			idx++
			if err != nil {
				return
			}

			size := 1
			if c == '\t' {
				size = 8 - (lineLength+8)%8
			}
			lineLength += size
			l.sum += 1

			l.lines[cline].content = append(l.lines[cline].content, Char{
				c:    c,
				size: size,
			})

			if c == '\n' {
				break
			}
		}
		cont = lineLength >= width
	}
}

func (l *LineView) Lines() []Line {
	return l.lines
}

func (l *LineView) CursorRelView(c Cursor) int {
	if c.Begin < l.disp {
		return -1
	}
	if c.Begin >= l.disp+l.sum {
		return 1
	}
	return 0
}

func (l *LineView) CursorLeft(c Cursor) Cursor {
	if c.Begin == 0 {
		return c
	}

	c.Begin--
	c.End = c.Begin + 1

	for {
		rel := l.CursorRelView(c)
		if rel == 0 {
			break
		}
		if rel == -1 {
			l.ScrollBack(1)
		} else {
			l.ScrollForward(1)
		}
	}

	return c
}

func (l *LineView) CursorRight(c Cursor) Cursor {
	if c.Begin == l.buffer.Size()-1 {
		return c
	}

	c.Begin++
	c.End = c.Begin + 1

	for {
		rel := l.CursorRelView(c)
		if rel == 0 {
			break
		}
		if rel == -1 {
			l.ScrollBack(1)
		} else {
			l.ScrollForward(1)
		}
	}

	return c
}

func (l *LineView) CursorLine(c Cursor) (int, int, error) {
	if l.CursorRelView(c) != 0 {
		return 0, 0, errors.New("cursor not in view")
	}

	line := 0
	disp := l.disp
	for {
		s := len(l.lines[line].content)
		if disp+s > c.Begin {
			break
		}
		disp += s
		line++
	}
	return line, c.Begin - disp, nil
}

func (l *LineView) LineDisp(line int) (d int, e error) {
	if line > len(l.lines) {
		return 0, errors.New("no such line")
	}
	d = l.disp
	for i := range line {
		d += len(l.lines[i].content)
	}
	return
}

func (l *LineView) CursorUp(c Cursor) Cursor {
	line, ld, err := l.CursorLine(c)
	if err != nil {
		return c
	}

	if line == 0 {
		err := l.scrollBack()
		if err != nil {
			return c
		}
		l.Update(l.height, l.width)
		line++
	}

	line--
	linedisp, _ := l.LineDisp(line)
	if ld >= len(l.lines[line].content) {
		c.Begin = linedisp + len(l.lines[line].content) - 1
	} else {
		c.Begin = linedisp + ld
	}
	c.End = c.Begin + 1

	return c
}

func (l *LineView) CursorDown(c Cursor) Cursor {
	line, ld, err := l.CursorLine(c)
	if err != nil {
		return c
	}

	if line == len(l.lines)-1 {
		if l.disp+l.Size() == l.buffer.Size() {
			return c
		}
		err := l.scrollForward()
		if err != nil {
			return c
		}
		l.Update(l.height, l.width)
		line--
	}

	line++
	linedisp, _ := l.LineDisp(line)
	if ld >= len(l.lines[line].content) {
		c.Begin = linedisp + len(l.lines[line].content) - 1
	} else {
		c.Begin = linedisp + ld
	}
	c.End = c.Begin + 1

	return c
}

func (l *LineView) scrollBack() error {
	if l.firstLine == 0 {
		return errors.New("cannot scroll back on the first line")
	}
	l.firstLine--
	disp := l.disp - 1
	for {
		if disp-1 < 0 {
			break
		}
		c, err := l.buffer.Get(disp - 1)
		if err != nil {
			break
		}
		if c == '\n' {
			break
		}
		disp--
	}
	l.disp = disp
	return nil
}

func (l *LineView) ScrollBack(amount int) {
	for range amount {
		if l.scrollBack() != nil {
			return
		}
	}
	l.Update(l.height, l.width)
}

func (l *LineView) scrollForward() error {
	disp := l.disp
	for {
		c, err := l.buffer.Get(disp)
		// Don't scroll.
		if err != nil {
			return errors.New("cannot scroll forward on the last line")
		}
		disp++
		if c == '\n' {
			_, err := l.buffer.Get(disp)
			// Don't scroll.
			if err != nil {
				return errors.New("cannot scroll forward on the last line")
			}
			break
		}
	}
	l.firstLine++
	l.disp = disp
	return nil
}

func (l *LineView) ScrollForward(amount int) {
	for range amount {
		if l.scrollForward() != nil {
			return
		}
	}
	l.Update(l.height, l.width)
}
