package main

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"slices"
)

type editor struct {
	scr     tcell.Screen
	colors  colorscheme
	minibuf minibuffer
	top     topBar
	buffers []*buffer
	window  *window
}

type size struct {
	width  int
	height int
}

type position struct {
	x int
	y int
}

func (e *editor) render() {
	e.drawBars()
	e.drawWindows()
	e.drawMinibuffer()
}

func initEditorOrPanic() *editor {
	scr, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := scr.Init(); err != nil {
		panic(err)
	}

	scr.EnableMouse()
	scr.EnablePaste()

	scr.SetTitle("ḥăzā")

	scr.SetStyle(tcell.StyleDefault)

	scr.Clear()

	e := editor{
		scr: scr,
	}

	e.loadDefaultColorscheme()
	e.initFirstWindow()

	return &e
}

func (e *editor) initScratchBuffer() {
	buf, _ := e.addBuffer(`Welcome to the ḥăzā' text editor.
This is a scratch buffer you can use for text that will not be saved.
`, "scratch/")
	e.window.buf = buf
}

func (e *editor) initFirstWindow() {
	e.window = new(window)
	e.recalculateWindowsPos()
	e.window.firstLine = 1
	e.initScratchBuffer()
	e.window.buf, _ = e.getBufferByName("scratch/")
}

func (e *editor) recalculateWindowsPos() {
	maxWidth, maxHeight := e.scr.Size()
	maxHeight -= 2

	e.window.pos = position{0, 1}
	e.window.sz = size{maxWidth, maxHeight}
}

func (e *editor) mainLoop() {
	defer func() {
		maybePanic := recover()
		e.scr.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}()

	for {
		e.render()
		e.scr.Show()

		e.defaultEventHandler(e.scr.PollEvent())
	}
}

func (e *editor) defaultEventHandler(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		e.recalculateWindowsPos()
		e.render()
		e.scr.Sync()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlQ:
			e.quit()
		case tcell.KeyCtrlW:
			e.commandDelete()
		case tcell.KeyCtrlG:
			e.clearMinibuf()
		case tcell.KeyCtrlP:
			cmd, cancelled := e.commandMode(":")
			if !cancelled {
				e.executeCommand(cmd)
			}
		}
	}
}

func (e *editor) executeCommand(command string) {
	switch command[:len(command)-1] {
	case "quit":
		e.quit()
	case "new":
		e.commandNew()
	case "delete":
		e.commandDelete()
	}
}

func (e *editor) commandNew() {
	buf := e.addUnnamedBuffer()
	e.gotoBuffer(e.window, buf)
}

func (e *editor) commandDelete() {
	buf := e.window.buf

	// Don't delete scratch.
	if buf == e.buffers[0] {
		e.warn("Refusing to delete scratch buffer.")
		return
	}

	idx := 0
	for i, b := range e.buffers {
		if b == buf {
			idx = i
			break
		}
	}

	e.buffers = slices.Delete(e.buffers, idx, idx+1)

	e.gotoBuffer(e.window, e.buffers[idx-1])
}

func (e *editor) quit() {
	e.scr.Fini()
	os.Exit(0)
}
