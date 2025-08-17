package main

type Char struct {
	c    rune
	size int
}

type Line struct {
	content      []Char
	continuation bool
}

type LineView struct {
	buffer    Buffer
	firstLine int
	disp      int
	lines     []Line
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

func (l *LineView) Update(height, width int) {
	l.lines = make([]Line, 0, height)
	idx := l.disp
	cline := -1
	cont := false

	for cline < height {
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
