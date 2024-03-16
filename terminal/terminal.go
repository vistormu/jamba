package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

type Terminal struct {
    Original Termios
    Modified *Termios
    NCols int
    NRows int
}

type winsize struct {
    Row, Col uint16
}

// ======
// PUBLIC
// ======
func (t *Terminal) Restore() error {
    fd := os.Stdout.Fd()
    _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(&t.Original)))
    if errno != 0 {
        return errno
    }

    return nil
}

func New() (*Terminal, error) {
    t := &Terminal{}

    fd := os.Stdout.Fd()
    termios, err := getTermios(fd)
    if err != nil {
        return nil, err
    }
    
    t.Original = *termios
    t.Modified = termios

    t.enableRawMode()
    err = t.getWindowSize()
    if err != nil {
        return nil, err
    }

    err = setTermios(fd, t.Modified)
    if err != nil {
        return nil, err
    }

    return t, nil
}

// =======
// PRIVATE
// =======
func (t *Terminal) enableRawMode() {
    t.Modified.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
    t.Modified.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
    t.Modified.Cflag |= syscall.CS8
    t.Modified.Oflag &^= syscall.OPOST
    t.Modified.Cc[syscall.VMIN+1] = 0
    t.Modified.Cc[syscall.VTIME+1] = 1
}
func (t *Terminal) getWindowSize() error {
    ws := &winsize{}
    _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(ws)))
    if errno != 0 {
        return errno
    }

    t.NCols = int(ws.Col)
    t.NRows = int(ws.Row)

    return nil
}
