package key

import (
    "fmt"
)

type Key int
const (
    Tab Key = 9
    Enter Key = 13
    Escape Key = 27
    Backspace Key = 127
    ArrowUp Key = 1000 + iota
    ArrowDown
    ArrowRight
    ArrowLeft
)

var keyNames = map[Key]string{
    Tab: "<Tab>",
    Enter: "<Enter>",
    Escape: "<Esc>" ,
    Backspace: "<BS>",
    ArrowUp: "<Up>", 
    ArrowDown: "<Down>",
    ArrowRight: "<Right>",
    ArrowLeft: "<Left>",
}

const (
    EscArrowUp = "\x1b[A"
    EscArrowDown = "\x1b[B"
    EscArrowRight = "\x1b[C"
    EscArrowLeft = "\x1b[D"
)

func New(c interface{}) (Key, error) {
    switch c := c.(type) {
    case byte:
        return Key(c), nil
    case [3]byte:
        escSeq := string(c[:])
        switch escSeq {
        case EscArrowUp:
            return ArrowUp, nil
        case EscArrowDown:
            return ArrowDown, nil
        case EscArrowRight:
            return ArrowRight, nil
        case EscArrowLeft:
            return ArrowLeft, nil
        default:
            if c[0] == 27 {
                return Escape, nil
            }
        }
    }

    return 0, fmt.Errorf("Unknown key type") 
}

func (k Key) String() string {
    name, ok := keyNames[k]
    if !ok {
        return fmt.Sprintf("%c", k)
    }
    return name
}
