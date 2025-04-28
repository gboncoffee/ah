package main

type topBarButton struct {
	buf  *buffer
	maxx int
}

type topBar struct {
	buttons []topBarButton
}

func (e *editor) drawBars() {
	width, height := e.scr.Size()

	e.drawTopBar(width)

	bh := height - 2
	for i := range width {
		e.scr.SetContent(i, bh, ' ', nil, e.colors.bar)
	}
}

func (e *editor) drawTopBar(width int) {
	e.top.buttons = []topBarButton{}

	// For simplicity we clear the bar before.
	for i := range width {
		e.scr.SetContent(i, 0, ' ', nil, e.colors.bar)
	}

	newx := 1
	for _, b := range e.buffers {
		oldx := newx
		name := b.name

		if b == e.window.buf {
			name = "[" + name + "]"
			newx += len(name)
			oldx--
		} else {
			newx += len(name) + 2
		}

		if newx >= width {
			break
		}

		e.top.buttons = append(e.top.buttons, topBarButton{
			maxx: newx,
			buf:  b,
		})

		var i int
		for _, ch := range name {
			e.scr.SetContent(oldx+i, 0, ch, nil, e.colors.bar)
			i++
		}
	}
}
