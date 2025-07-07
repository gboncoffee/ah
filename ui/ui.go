package ui

import "github.com/gdamore/tcell/v2"

type Ui struct {
	screen tcell.Screen
	state  State

	updates chan Update
	exit    chan struct{}

	defaultStyle tcell.Style
}

type mode int

const (
	modeEditor = iota
	modeMinibuffer
	modeSide
)

type Update func(state *State)

func New() *Ui {
	scr, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	if err := scr.Init(); err != nil {
		panic(err)
	}

	ui := new(Ui)

	ui.screen = scr
	ui.defaultStyle = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	ui.screen.SetStyle(ui.defaultStyle)

	ui.updates = make(chan Update)
	ui.exit = make(chan struct{})

	return ui
}

func (ui *Ui) Start(start func(ui *State)) {
	go ui.Update(start)
	ui.mainloop()
}

func (ui *Ui) traceMeHarder(finish bool) {
	r := recover()
	if r != nil || finish {
		ui.screen.Fini()
	}

	if r != nil {
		panic(r)
	}
}

func (ui *Ui) mainloop() {
	defer ui.traceMeHarder(true)

	input := make(chan tcell.Event)
	go ui.screen.ChannelEvents(input, make(chan struct{}))

	ui.render()
	for {
		select {
		case update := <-ui.updates:
			update(&ui.state)
			ui.render()
		case <-ui.exit:
			return
		case ev := <-input:
			go ui.input(EventFromTcell(ev))
		}
	}
}

func (ui *Ui) Update(f Update) {
	ui.updates <- f
}

func (ui *Ui) Exit() {
	close(ui.exit)
}

func (ui *Ui) DefaultStyle(style tcell.Style) {
	ui.defaultStyle = style
	ui.screen.SetStyle(ui.defaultStyle)
}
