package pty

import "os"

type Terminal struct {
	Pty      *os.File
	Tty      *os.File
	Recorder []*os.File
	Log      *os.File
	oldState State
}

func (t *Terminal) Write(b []byte) (n int, err error) {
	return t.Pty.Write(b)
}

func (t *Terminal) Read(b []byte) (n int, err error) {
	n, err = t.Pty.Read(b)
	for _, r := range t.Recorder {
		if n, err := r.Write(b); err != nil {
			return n, err
		}
	}
	return n, err
}

func (t *Terminal) SetWinSize(x, y int) error {
	return SetWinSize(t.Pty, uint16(x), uint16(y))
}

func (t *Terminal) GetWinSize() (x, y int, err error) {
	return GetWinSize(t.Pty)
}

func (t *Terminal) ResetWinSize() error {
	if x, y, err := t.GetWinSize(); err != nil {
		return err
	} else {
		return t.SetWinSize(x, y)
	}
}

func (t *Terminal) SetRaw() (err error) {
	if oldState, err := Tcgetattr(t.Pty); err != nil {
		return err
	} else {
		t.oldState = *oldState
	}
	return SetRaw(t.Pty)
}

func (t *Terminal) SetCBreak() (err error) {
	if oldState, err := Tcgetattr(t.Pty); err != nil {
		return err
	} else {
		t.oldState = *oldState
	}
	return SetCBreak(t.Pty)
}

func (t *Terminal) Restore() (err error) {
	return Tcsetattr(t.Pty, &t.oldState)
}

func (t *Terminal) SendIntr() (err error) {
	c, err := GetControlChar(t.Pty, "CTRL-C")
	_, err = t.Write([]byte{c})
	return
}

func (t *Terminal) SendEOF() (err error) {
	c, err := GetControlChar(t.Pty, "EOF")
	_, err = t.Write([]byte{c})
	return
}

func (t *Terminal) Close() (err error) {
	stdout.Reset()
	for _, r := range t.Recorder {
		r.Close()
	}
	err = t.Tty.Close()
	err = t.Pty.Close()

	return
}

func NewTerminal() (term *Terminal, err error) {
	pty, tty, err := Open()
	term = &Terminal{Pty: pty, Tty: tty}
	return
}
