package main

import "github.com/gdamore/tcell/v2"

type colorscheme struct {
	def        tcell.Style
	bar        tcell.Style
	minibuffer tcell.Style
	lineNums   tcell.Style
	eob        tcell.Style
}

func (e *editor) loadDefaultColorscheme() {
	c := &e.colors

	c.def = tcell.StyleDefault.
		Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	c.bar = tcell.StyleDefault.
		Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	c.minibuffer = tcell.StyleDefault.
		Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	c.lineNums = tcell.StyleDefault.
		Background(tcell.ColorGray).Foreground(tcell.ColorYellow)

	c.eob = tcell.StyleDefault.
		Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
}
