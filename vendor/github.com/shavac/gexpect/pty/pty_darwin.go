package pty

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

const (
	_IOC_VOID    uintptr = 0x20000000
	_IOC_OUT     uintptr = 0x40000000
	_IOC_IN      uintptr = 0x80000000
	_IOC_IN_OUT  uintptr = _IOC_OUT | _IOC_IN
	_IOC_DIRMASK         = _IOC_VOID | _IOC_OUT | _IOC_IN

	_IOC_PARAM_SHIFT = 13
	_IOC_PARAM_MASK  = (1 << _IOC_PARAM_SHIFT) - 1
)

func Open() (pty, tty *os.File, err error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	err = grantpt(p)
	if err != nil {
		return nil, nil, err
	}

	err = unlockpt(p)
	if err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	return p, t, nil
}

func ptsname(f *os.File) (string, error) {
	n := make([]byte, _IOC_PARM_LEN(syscall.TIOCPTYGNAME))

	err := ioctl(f.Fd(), syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&n[0])))
	if err != nil {
		return "", err
	}

	for i, c := range n {
		if c == 0 {
			return string(n[:i]), nil
		}
	}
	return "", errors.New("TIOCPTYGNAME string not NUL-terminated")
}

func grantpt(f *os.File) error {
	return ioctl(f.Fd(), syscall.TIOCPTYGRANT, 0)
}

func unlockpt(f *os.File) error {
	return ioctl(f.Fd(), syscall.TIOCPTYUNLK, 0)
}

func _IOC_PARM_LEN(ioctl uintptr) uintptr {
	return (ioctl >> 16) & _IOC_PARAM_MASK
}

func _IOC(inout uintptr, group byte, ioctl_num uintptr, param_len uintptr) uintptr {
	return inout | (param_len&_IOC_PARAM_MASK)<<16 | uintptr(group)<<8 | ioctl_num
}

func _IO(group byte, ioctl_num uintptr) uintptr {
	return _IOC(_IOC_VOID, group, ioctl_num, 0)
}

func _IOR(group byte, ioctl_num uintptr, param_len uintptr) uintptr {
	return _IOC(_IOC_OUT, group, ioctl_num, param_len)
}

func _IOW(group byte, ioctl_num uintptr, param_len uintptr) uintptr {
	return _IOC(_IOC_IN, group, ioctl_num, param_len)
}

func _IOWR(group byte, ioctl_num uintptr, param_len uintptr) uintptr {
	return _IOC(_IOC_IN_OUT, group, ioctl_num, param_len)
}

func ioctl(fd, cmd, ptr uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if e != 0 {
		return e
	}
	return nil
}
