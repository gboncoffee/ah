package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

type minibuffer struct {
	prompt  string
	content string
	cursor  int
	active  bool
	warning bool
}

func (e *editor) drawMinibuffer() {
	width, height := e.scr.Size()

	for i := range width {
		e.scr.SetContent(i, height-1, ' ', nil, e.colors.minibuffer)
	}

	if e.minibuf.active {
		var i int
		for _, ch := range e.minibuf.prompt {
			e.scr.SetContent(i, height-1, ch, nil, e.colors.minibuffer)
			i++
			if i >= width {
				break
			}
		}

		var j int
		for _, ch := range e.minibuf.content {
			color := e.colors.minibuffer
			if j == e.minibuf.cursor {
				color = color.Reverse(true)
			}
			e.scr.SetContent(j+i, height-1, ch, nil, color)
			j++
		}
	}
}

func (e *editor) commandMode(prompt string) (result string, cancelled bool) {
	e.minibuf.active = true
	e.minibuf.prompt = prompt
	e.minibuf.cursor = 0
	e.minibuf.content = "\n"

	for {
		e.render()
		e.scr.Show()

		switch ev := e.scr.PollEvent().(type) {
		case *tcell.EventResize:
			e.defaultEventHandler(ev)
		case *tcell.EventKey:
			quit, cancelled := e.commandModeKeyTyped(ev)
			if quit {
				e.minibuf.active = false
				return e.minibuf.content, cancelled
			}
		}
	}
}

func (e *editor) commandModeKeyTyped(
	ev *tcell.EventKey,
) (quit bool, cancelled bool) {

	switch ev.Key() {

	case tcell.KeyLeft:
		e.commandModeLeft()
		return false, false

	case tcell.KeyRight:
		e.commandModeRight()
		return false, false

	case tcell.KeyRune:
		e.commandModeInsert(ev.Rune())
		return false, false

	case tcell.KeyCtrlW:
		return true, true

	case tcell.KeyEnter:
		return true, false

	default:
		return false, false
	}
}

func (e *editor) commandModeLeft() {
	if e.minibuf.cursor > 0 {
		e.minibuf.cursor--
	}
}

func (e *editor) commandModeRight() {
	if e.minibuf.cursor < len(e.minibuf.content)-1 {
		e.minibuf.cursor++
	}
}

func (e *editor) commandModeInsert(ch rune) {
	var newContent strings.Builder

	var i int
	for _, oldch := range e.minibuf.content {
		if i == e.minibuf.cursor {
			newContent.WriteRune(ch)
		}
		newContent.WriteRune(oldch)
		i++
	}

	e.minibuf.content = newContent.String()
	e.minibuf.cursor++
}
