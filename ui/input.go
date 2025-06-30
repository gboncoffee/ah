package ui

import "github.com/gdamore/tcell/v2"

type Key tcell.Key

const (
	KeyUp        Key = Key(tcell.KeyUp)
	KeyDown      Key = Key(tcell.KeyDown)
	KeyRight     Key = Key(tcell.KeyRight)
	KeyLeft      Key = Key(tcell.KeyLeft)
	KeyUpLeft    Key = Key(tcell.KeyUpLeft)
	KeyUpRight   Key = Key(tcell.KeyUpRight)
	KeyDownLeft  Key = Key(tcell.KeyDownLeft)
	KeyDownRight Key = Key(tcell.KeyDownRight)
	KeyCenter    Key = Key(tcell.KeyCenter)
	KeyPgUp      Key = Key(tcell.KeyPgUp)
	KeyPgDn      Key = Key(tcell.KeyPgDn)
	KeyHome      Key = Key(tcell.KeyHome)
	KeyEnd       Key = Key(tcell.KeyEnd)
	KeyInsert    Key = Key(tcell.KeyInsert)
	KeyDelete    Key = Key(tcell.KeyDelete)
	KeyHelp      Key = Key(tcell.KeyHelp)
	KeyExit      Key = Key(tcell.KeyExit)
	KeyClear     Key = Key(tcell.KeyClear)
	KeyCancel    Key = Key(tcell.KeyCancel)
	KeyPrint     Key = Key(tcell.KeyPrint)
	KeyPause     Key = Key(tcell.KeyPause)
	KeyBacktab   Key = Key(tcell.KeyBacktab)
	KeyF1        Key = Key(tcell.KeyF1)
	KeyF2        Key = Key(tcell.KeyF2)
	KeyF3        Key = Key(tcell.KeyF3)
	KeyF4        Key = Key(tcell.KeyF4)
	KeyF5        Key = Key(tcell.KeyF5)
	KeyF6        Key = Key(tcell.KeyF6)
	KeyF7        Key = Key(tcell.KeyF7)
	KeyF8        Key = Key(tcell.KeyF8)
	KeyF9        Key = Key(tcell.KeyF9)
	KeyF10       Key = Key(tcell.KeyF10)
	KeyF11       Key = Key(tcell.KeyF11)
	KeyF12       Key = Key(tcell.KeyF12)
	// No I won't support more than 12 function keys lol.

	KeyCtrlSpace      Key = Key(tcell.KeyCtrlSpace)
	KeyCtrlA          Key = Key(tcell.KeyCtrlA)
	KeyCtrlB          Key = Key(tcell.KeyCtrlB)
	KeyCtrlC          Key = Key(tcell.KeyCtrlC)
	KeyCtrlD          Key = Key(tcell.KeyCtrlD)
	KeyCtrlE          Key = Key(tcell.KeyCtrlE)
	KeyCtrlF          Key = Key(tcell.KeyCtrlF)
	KeyCtrlG          Key = Key(tcell.KeyCtrlG)
	KeyCtrlH          Key = Key(tcell.KeyCtrlH)
	KeyCtrlI          Key = Key(tcell.KeyCtrlI)
	KeyCtrlJ          Key = Key(tcell.KeyCtrlJ)
	KeyCtrlK          Key = Key(tcell.KeyCtrlK)
	KeyCtrlL          Key = Key(tcell.KeyCtrlL)
	KeyCtrlM          Key = Key(tcell.KeyCtrlM)
	KeyCtrlN          Key = Key(tcell.KeyCtrlN)
	KeyCtrlO          Key = Key(tcell.KeyCtrlO)
	KeyCtrlP          Key = Key(tcell.KeyCtrlP)
	KeyCtrlQ          Key = Key(tcell.KeyCtrlQ)
	KeyCtrlR          Key = Key(tcell.KeyCtrlR)
	KeyCtrlS          Key = Key(tcell.KeyCtrlS)
	KeyCtrlT          Key = Key(tcell.KeyCtrlT)
	KeyCtrlU          Key = Key(tcell.KeyCtrlU)
	KeyCtrlV          Key = Key(tcell.KeyCtrlV)
	KeyCtrlW          Key = Key(tcell.KeyCtrlW)
	KeyCtrlX          Key = Key(tcell.KeyCtrlX)
	KeyCtrlY          Key = Key(tcell.KeyCtrlY)
	KeyCtrlZ          Key = Key(tcell.KeyCtrlZ)
	KeyCtrlLeftSq     Key = Key(tcell.KeyCtrlLeftSq)
	KeyCtrlBackslash  Key = Key(tcell.KeyCtrlBackslash)
	KeyCtrlRightSq    Key = Key(tcell.KeyCtrlRightSq)
	KeyCtrlCarat      Key = Key(tcell.KeyCtrlCarat)
	KeyCtrlUnderscore Key = Key(tcell.KeyCtrlUnderscore)
)

type Event interface {
	Underlying() tcell.Event
}

type MouseEvent interface {
	Position() (x, y int)
}

type MouseRightClick struct {
	x, y       int
	underlying *tcell.EventMouse
}

func (m *MouseRightClick) Underlying() tcell.Event {
	return m.underlying
}

func (m *MouseRightClick) Position() (int, int) {
	return m.x, m.y
}

type MouseLeftClick struct {
	x, y       int
	underlying *tcell.EventMouse
}

func (m *MouseLeftClick) Underlying() tcell.Event {
	return m.underlying
}

func (m *MouseLeftClick) Position() (int, int) {
	return m.x, m.y
}

type MouseMiddleClick struct {
	x, y       int
	underlying *tcell.EventMouse
}

func (m *MouseMiddleClick) Underlying() tcell.Event {
	return m.underlying
}

func (m *MouseMiddleClick) Position() (int, int) {
	return m.x, m.y
}

type MouseScrollUp struct {
	x, y       int
	underlying *tcell.EventMouse
}

func (m *MouseScrollUp) Underlying() tcell.Event {
	return m.underlying
}

func (m *MouseScrollUp) Position() (int, int) {
	return m.x, m.y
}

type MouseScrollDown struct {
	x, y       int
	underlying *tcell.EventMouse
}

func (m *MouseScrollDown) Underlying() tcell.Event {
	return m.underlying
}

func (m *MouseScrollDown) Position() (int, int) {
	return m.x, m.y
}

type MouseScrollLeft struct {
	x, y       int
	underlying *tcell.EventMouse
}

func (m *MouseScrollLeft) Underlying() tcell.Event {
	return m.underlying
}

func (m *MouseScrollLeft) Position() (int, int) {
	return m.x, m.y
}

type MouseScrollRight struct {
	x, y       int
	underlying *tcell.EventMouse
}

func (m *MouseScrollRight) Underlying() tcell.Event {
	return m.underlying
}

func (m *MouseScrollRight) Position() (int, int) {
	return m.x, m.y
}

type KeyPress struct {
	Key        Key
	underlying *tcell.EventKey
}

func (k *KeyPress) Underlying() tcell.Event {
	return k.underlying
}

type RuneEntered struct {
	Rune       rune
	underlying *tcell.EventKey
}

func (r *RuneEntered) Underlying() tcell.Event {
	return r.underlying
}

type Resize struct {
	Width, Height int
	underlying    *tcell.EventResize
}

func (r *Resize) Underlying() tcell.Event {
	return r.underlying
}

func EventFromTcell(ev tcell.Event) Event {
	switch e := ev.(type) {
	case *tcell.EventKey:
		return keypressFromTcell(e)
	case *tcell.EventResize:
		return resizeFromTcell(e)
	case *tcell.EventMouse:
		return mouseFromTcell(e)
	}

	return nil
}

func keypressFromTcell(e *tcell.EventKey) Event {
	key := e.Key()
	if key == tcell.KeyRune {
		return &RuneEntered{
			underlying: e,
			Rune:       e.Rune(),
		}
	}
	return &KeyPress{
		underlying: e,
		Key:        Key(key),
	}
}

func resizeFromTcell(e *tcell.EventResize) Event {
	w, h := e.Size()
	return &Resize{
		Width:      w,
		Height:     h,
		underlying: e,
	}
}

func mouseFromTcell(ev *tcell.EventMouse) Event {
	mask := ev.Buttons()
	x, y := ev.Position()
	if mask&tcell.ButtonSecondary != 0 {
		return &MouseRightClick{
			underlying: ev,
			x:          x,
			y:          y,
		}
	}
	if mask&tcell.ButtonMiddle != 0 {
		return &MouseRightClick{
			underlying: ev,
			x:          x,
			y:          y,
		}
	}
	if mask&tcell.WheelUp != 0 {
		return &MouseScrollUp{
			underlying: ev,
			x:          x,
			y:          y,
		}
	}
	if mask&tcell.WheelDown != 0 {
		return &MouseScrollDown{
			underlying: ev,
			x:          x,
			y:          y,
		}
	}
	if mask&tcell.WheelLeft != 0 {
		return &MouseScrollLeft{
			underlying: ev,
			x:          x,
			y:          y,
		}
	}
	if mask&tcell.WheelRight != 0 {
		return &MouseScrollRight{
			underlying: ev,
			x:          x,
			y:          y,
		}
	}

	return &MouseLeftClick{
		underlying: ev,
		x:          x,
		y:          y,
	}
}
