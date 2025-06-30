package ui

type State struct {
	Editor     EditorView
	Minibuffer EditorView
	mode       mode
}

func (ui *Ui) render() {
	ui.screen.Clear()
	w, h := ui.screen.Size()
	if ui.state.Editor != nil {
		ui.renderEditor(ui.state.Editor, 0, 0, w, h-1, ui.state.mode == modeEditor)
	}
	if ui.state.Minibuffer != nil {
		ui.renderEditor(ui.state.Minibuffer, 0, h-1, w, 1, ui.state.mode == modeMinibuffer)
	}
	ui.screen.Show()
}
