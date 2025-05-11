package main

import (
	"bufio"
	"io"
	"math"
	"strconv"
)

type window struct {
	pos          position
	sz           size
	buf          *buffer
	firstLine    int
	lineOffset   int
	cursorOffset int
}

func (e *editor) drawWindows() {
	e.drawWindow(e.window)
}

func (e *editor) drawWindow(w *window) {
	reader := bufio.NewReader(w.buf.content.OffsetReader(w.lineOffset))

	sizeOfLineNumbers := max(int(
		math.Floor(math.Log10(float64(w.buf.numberOfLines))),
	), 3)

	var end int
	eof := false
	offsetInBuffer := w.lineOffset
	for i := range w.sz.height - 1 {
		end = i
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			panic(err)
		}

		lineNum := strconv.Itoa(w.firstLine + i)
		pad := sizeOfLineNumbers - len(lineNum)
		var padEnd int
		for j := range pad {
			e.scr.SetContent(w.pos.x+j, w.pos.y+i, ' ', nil, e.colors.lineNums)
			padEnd = j
		}
		// Safe because it's always ASCII.
		for j, ch := range lineNum {
			e.scr.SetContent(
				w.pos.x+padEnd+j+1, w.pos.y+i, ch, nil, e.colors.lineNums,
			)
		}

		var end int
		j := 0

		e.scr.SetContent(
			w.pos.x+j+sizeOfLineNumbers,
			w.pos.y+i, ' ', nil, e.colors.lineNums,
		)

		for _, ch := range line {
			end = j
			if j >= w.sz.width-sizeOfLineNumbers-1 {
				break
			}

			color := e.colors.def
			if !e.minibuf.active && e.window == w &&
				offsetInBuffer == w.cursorOffset {
				color = color.Reverse(true)
			}

			e.scr.SetContent(
				w.pos.x+j+sizeOfLineNumbers+1,
				w.pos.y+i, ch, nil, color,
			)
			j++
			offsetInBuffer++
		}
		for j := end + 1; j < w.sz.width-sizeOfLineNumbers-1; j++ {
			e.scr.SetContent(
				w.pos.x+j+sizeOfLineNumbers+1,
				w.pos.y+i, ' ', nil, e.colors.def,
			)
		}

		if eof {
			break
		}

		// Account for the line ending.
		offsetInBuffer += 1
	}

	for end < w.sz.height-1 {
		e.scr.SetContent(w.pos.x, w.pos.y+end, '~', nil, e.colors.eob)
		for i := 1; i <= sizeOfLineNumbers; i++ {
			e.scr.SetContent(w.pos.x+i, w.pos.y+end, ' ', nil, e.colors.eob)
		}
		for i := sizeOfLineNumbers + 1; i < w.sz.width; i++ {
			e.scr.SetContent(w.pos.x+i, w.pos.y+end, ' ', nil, e.colors.def)
		}
		end++
	}
}

func (e *editor) gotoBuffer(win *window, buf *buffer) {
	e.window.buf = buf
	e.window.firstLine = 1
	e.window.cursorOffset = 0
	e.window.lineOffset = 0
}

func (w *window) runePress(r rune) {
	w.buf.insert(w.cursorOffset, r)
}
