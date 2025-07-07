package ui

func (ui *Ui) input(ev Event) {
	switch e := ev.(type) {
	case *KeyPress:
		ui.keypress(e)
	case *RuneEntered:
		ui.runeEntered(e)
	case *Resize:
		ui.resize(e)
	case *MouseRightClick:
		ui.rightClick(e)
	case *MouseLeftClick:
		ui.leftClick(e)
	case *MouseMiddleClick:
		ui.middleClick(e)
	case *MouseScrollUp:
		ui.scrollUp(e)
	case *MouseScrollDown:
		ui.scrollDown(e)
	case *MouseScrollLeft:
		ui.scrollLeft(e)
	case *MouseScrollRight:
		ui.scrollRight(e)
	}
}

func (ui *Ui) keypress(e *KeyPress) {
	if e.Key == KeyCtrlC {
		ui.Exit()
	}

	if ui.state.mode == modeEditor && ui.state.Editor != nil {
		ui.state.Editor.Event(e)
	}
}

func (ui *Ui) runeEntered(e *RuneEntered) {
	if ui.state.mode == modeEditor && ui.state.Editor != nil {
		ui.state.Editor.Event(e)
	}
}

func (ui *Ui) resize(e *Resize) {}

func (ui *Ui) rightClick(e *MouseRightClick) {}

func (ui *Ui) leftClick(e *MouseLeftClick) {}

func (ui *Ui) middleClick(e *MouseMiddleClick) {}

func (ui *Ui) scrollUp(e *MouseScrollUp) {}

func (ui *Ui) scrollDown(e *MouseScrollDown) {}

func (ui *Ui) scrollLeft(e *MouseScrollLeft) {}

func (ui *Ui) scrollRight(e *MouseScrollRight) {}
