package terminal

import (
    "syscall"
    "unsafe"
)

type Termios struct {
    Iflag, Oflag, Cflag, Lflag uint32
    Cc [20]byte
    Ispeed, Ospeed uint32
}

func getTermios(fd uintptr) (*Termios, error) {
    termios := &Termios{}
    _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)))
    if errno != 0 {
        return nil, errno
    }

    return termios, nil
}

func setTermios(fd uintptr, termios *Termios) error {
    _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS+1), uintptr(unsafe.Pointer(termios)))
    if errno != 0 {
        return errno
    }

    return nil
}
