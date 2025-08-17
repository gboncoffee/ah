package ui

import "github.com/gdamore/tcell/v2"

type State struct {
	Editor     EditorView
	Minibuffer EditorView
	Warning    string
	Message    string
	focus      any
	width      int
	height     int
}

func (s *State) EditorHeight() int {
	return s.height
}

func (ui *Ui) render() {
	ui.screen.Clear()
	w, h := ui.screen.Size()
	ui.state.width, ui.state.height = w, h
	if ui.state.Editor != nil {
		ui.renderEditor(
			ui.state.Editor,
			0, 0, w, h-1,
			ui.state.Editor == ui.state.focus,
		)
	}

	if ui.state.Warning != "" {
		ui.renderMessage(ui.state.Warning, ui.warningStyle, w, h)
		ui.state.Warning = ""
	} else if ui.state.Message != "" {
		ui.renderMessage(ui.state.Message, ui.messageStyle, w, h)
		ui.state.Message = ""
	} else if ui.state.Minibuffer != nil {
		ui.renderEditor(
			ui.state.Minibuffer,
			0, h-1, w, 1,
			ui.state.Minibuffer == ui.state.focus,
		)
	}
	ui.screen.Show()
}

func (ui *Ui) renderMessage(message string, style tcell.Style, w, h int) {
	i := 0
	for _, c := range message {
		if i >= w {
			break
		}

		ui.screen.SetContent(i, h-1, c, nil, style)
		i++
	}

	for i < w {
		ui.screen.SetContent(i, h-1, ' ', nil, style)
		i++
	}
}

func (s *State) Focus(thing any) {
	s.focus = thing
}
