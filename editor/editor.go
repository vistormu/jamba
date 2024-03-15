package editor

import (
	"fmt"
	"io"
	"jamba/key"
	"os"
	"strings"
)

type EditorMode int
const (
    Normal EditorMode = iota
    Insert
    Command
)
func (e EditorMode) String() string {
    return [...]string{"NORMAL", "INSERT", "COMMAND"}[e]
}

type Editor struct {
    // input
    Filename string
    Content []string

    // output
    Dirty bool
    Info string

    // cursor
    CursorX int
    CursorY int

    // mode
    Mode EditorMode

    // command
    Command string

    // internal
    buffer *key.Key
}

func New(filename string) (*Editor, error) {
    if filename == "" {
        return nil, fmt.Errorf("No filename provided")
    }
    
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(strings.TrimSpace(string(content)), "\n")
    
    return &Editor{
        Filename: filename,
        Content: lines,
        Mode: Normal,
    }, nil
}
func (e *Editor) Update(k key.Key) error {
    return e.processKeyPress(k)
}

// =======
// private
// =======
// state
func (e *Editor) save() error {
    file, err := os.Create(e.Filename)
    if err != nil {
        e.Info = "Error: " + err.Error()
        return err
    }
    defer file.Close()

    for _, line := range e.Content {
        file.WriteString(line + "\n")
    }

    e.Info = "File saved"
    e.Dirty = false

    return nil
}

// key press
func (e *Editor) processKeyPress(k key.Key) error {
    e.moveCursor(k)

    switch e.Mode {
    case Normal:
        return e.processNormalModeKeyPress(k)
    case Insert:
        return e.processInsertModeKeyPress(k)
    case Command:
        return e.processCommandModeKeyPress(k)
    }

    return fmt.Errorf("Unknown mode")
}

// cursor movement
func (e *Editor) moveCursor(k key.Key) {
    contentLen := len(e.Content)

    switch k {
    case key.ArrowUp:
        // content limit
        if e.CursorY == 0 {
            break
        }
        // vertical movement
        e.CursorY--
        // snap to end of line
        rowLen := len(e.Content[e.CursorY])
        prevRowLen := len(e.Content[e.CursorY + 1])
        if e.CursorX > rowLen || e.CursorX >= prevRowLen {
            e.CursorX = rowLen
        }

    case key.ArrowDown:
        // content limit
        if e.CursorY == contentLen - 1 {
            break
        }
        // vertical movement
        e.CursorY++
        // snap to end of line
        rowLen := len(e.Content[e.CursorY])
        prevRowLen := len(e.Content[e.CursorY - 1])
        if e.CursorX > rowLen || e.CursorX >= prevRowLen {
            e.CursorX = rowLen
        }

    case key.ArrowLeft:
        if e.CursorX > 0 {
            e.CursorX--
        }

    case key.ArrowRight:
        if e.CursorX < len(e.Content[e.CursorY]) {
            e.CursorX++
        }
    }
}

// normal mode
func (e *Editor) processNormalModeKeyPress(k key.Key) error {
    switch k {
    case 'i':
        e.Mode = Insert
    case ':':
        e.Mode = Command
    case 'o':
        e.CursorX = len(e.Content[e.CursorY])
        e.insertNewLine()
        e.Dirty = true
        e.CursorX = 0
        e.CursorY++
        e.Mode = Insert
    case 'A':
        e.CursorX = len(e.Content[e.CursorY])
        e.Mode = Insert
    case 'x':
        e.deleteChar()
    case 's':
        e.deleteChar()
        e.Mode = Insert
    case 'd':
        if e.buffer == nil {
            e.buffer = &k
            return nil
        }
        if *e.buffer == k {
            e.deleteLine()
            e.buffer = nil
        }
    case key.Escape:
        e.buffer = nil
    }

    return nil
} 

// insert mode
func (e *Editor) processInsertModeKeyPress(k key.Key) error {
    switch k {
    case key.Escape:
        e.Mode = Normal
    case key.Enter:
        e.insertNewLine()
        e.Dirty = true
        e.CursorX = 0
        e.CursorY++
    case key.Backspace:
        e.deleteChar()
    case key.Tab:
        for range 4 {
            e.insertChar(' ')
        }
    default:
        e.insertChar(k)
    }

    return nil
}
func (e *Editor) insertNewLine() {
    rowLen := len(e.Content[e.CursorY])

    // at the end of a line
    if e.CursorX >= rowLen {
        e.Content = append(e.Content[:e.CursorY+1], e.Content[e.CursorY:]...)
        e.Content[e.CursorY+1] = ""
        return
    }
    // in the middle of a line
    newLine := e.Content[e.CursorY][e.CursorX:]
    e.Content[e.CursorY] = e.Content[e.CursorY][:e.CursorX]
    e.Content = append(e.Content[:e.CursorY+1], e.Content[e.CursorY:]...)
    e.Content[e.CursorY+1] = newLine
}
func (e *Editor) deleteChar() {
    // at the beginning of the document
    if e.CursorX == 0  && e.CursorY == 0 {
        return
    }
    // normal backspace
    if e.CursorX > 0 {
        e.CursorX--
        e.Content[e.CursorY] = e.Content[e.CursorY][:e.CursorX] + e.Content[e.CursorY][e.CursorX+1:]
        e.Dirty = true
        return
    }
    // backspace at the beginning of a line
    if e.CursorX == 0 {
        e.CursorX = len(e.Content[e.CursorY-1])
        e.Content[e.CursorY-1] += e.Content[e.CursorY]
        e.Content = append(e.Content[:e.CursorY], e.Content[e.CursorY+1:]...)
        e.CursorY--
        e.Dirty = true
    }
}
func (e *Editor) insertChar(k key.Key) {
    if k < 32 || k > 126 {
        return
    }

    e.Content[e.CursorY] = e.Content[e.CursorY][:e.CursorX] + k.String() + e.Content[e.CursorY][e.CursorX:]
    e.CursorX++
    e.Dirty = true
}
func (e *Editor) deleteLine() {
    if len(e.Content) == 1 {
        e.Content[0] = ""
        e.CursorX = 0
        e.CursorY = 0
        e.Dirty = true
        return
    }
    if e.CursorY == 0 {
        e.Content = e.Content[1:]
        e.CursorX = len(e.Content[e.CursorY])
        return
    }
    if e.CursorY == len(e.Content) - 1 {
        e.Content = e.Content[:len(e.Content)-1]
        e.CursorY--
        e.CursorX = len(e.Content[e.CursorY])
        return
    }
    e.Content = append(e.Content[:e.CursorY], e.Content[e.CursorY+1:]...)
    e.CursorX = len(e.Content[e.CursorY])
    e.Dirty = true
}

// command mode
func (e *Editor) processCommandModeKeyPress(k key.Key) error {
    switch k {
    case key.Escape:
        e.Mode = Normal
        e.Command = ""
    case key.Enter:
        err := e.processCommand()
        e.Mode = Normal
        e.Command = ""
        return err
    case key.Backspace:
        if len(e.Command) > 0 {
            e.Command = e.Command[:len(e.Command)-1]
        }
    default:
        e.Command += k.String()
        e.Info = e.Command
    }

    return nil
}
func (e *Editor) processCommand() error {
    switch strings.ToLower(e.Command) {
    case "q":
        if e.Dirty {
            e.Info = "File has unsaved changes. Use :w to save or :q! to quit without saving"
            return nil
        }
        return io.EOF
    case "q!":
        return io.EOF
    case "w":
        return e.save()
    case "wq": 
        err := e.save()
        if err != nil {
            return err
        }
        return io.EOF
    }

    return nil
}
