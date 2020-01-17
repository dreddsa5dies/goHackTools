package pty

import (
	"os"
	"syscall"
	"unsafe"
	"errors"
)

const (
	IFLAG  = 0
	OFLAG  = 1
	CFLAG  = 2
	LFLAG  = 3
	ISPEED = 4
	OSPEED = 5
	CC     = 6
)

type ttySize struct {
	Rows   uint16
	Cols   uint16
	Xpixel uint16
	Ypixel uint16
}

type State struct {
	termios syscall.Termios
}

func SetWinSize(f *os.File, cols uint16, rows uint16) error {
	_, _, e := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&ttySize{rows, cols, 0, 0})),
		0, 0, 0,
	)
	if e != 0 {
		return syscall.ENOTTY
	}
	return nil
}

func GetWinSize(f *os.File) (width, height int, err error) {
	var dimensions ttySize
	if _, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dimensions)),
		0, 0, 0); err != 0 {
		return -1, -1, err
	}
	return int(dimensions.Cols), int(dimensions.Rows), nil
}

func SetRaw(f *os.File) (err error) {
	var state *State
	if state, err = Tcgetattr(f); err != nil {
		return
	}
	state.termios.Iflag &^= syscall.ISTRIP | syscall.ICRNL | syscall.IXON | syscall.BRKINT | syscall.INPCK
	state.termios.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	state.termios.Oflag &^= syscall.OPOST
	state.termios.Cflag &^= syscall.CSIZE | syscall.PARENB
	state.termios.Cflag |= syscall.CS8
	state.termios.Cc[syscall.VMIN] = 1
	state.termios.Cc[syscall.VTIME] = 0
	return Tcsetattr(f, state)
}

func SetCBreak(f *os.File) (err error) {
	var state *State
	if state, err = Tcgetattr(f); err != nil {
		return
	}
	state.termios.Lflag &^= syscall.ECHO | syscall.ICANON
	state.termios.Cc[syscall.VMIN] = 1
	state.termios.Cc[syscall.VTIME] = 0
	return Tcsetattr(f, state)
}

func Tcgetattr(f *os.File) (state *State, err error) {
	state = new(State)
	if _, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&state.termios)),
		0, 0, 0); err != 0 {
		return nil, err
	}
	return
}

func Tcsetattr(f *os.File, state *State) (err error) {
	if _, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(f.Fd()),
		syscall.TCSETS,
		uintptr(unsafe.Pointer(&state.termios)),
		0, 0, 0); err != 0 {
		return err
	}
	return
}

func GetControlChar(f *os.File, name string) (c byte, err error) {
	state, _ := Tcgetattr(f)
	switch name {
	case "EOF":
		c = state.termios.Cc[syscall.VEOF]
	case "CTRL-C":
		c = state.termios.Cc[syscall.VINTR]
	default:
		return 0, errors.New("No such controlling character.")
	}
	return byte(c), nil
}
