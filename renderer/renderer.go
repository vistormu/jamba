package renderer

import (
    "fmt"
    "jamba/editor"
)

const (
    clearScreen = "\x1b[2J"
    clearLine = "\x1b[K"
    cursorHide = "\x1b[?25l"
    cursorShow = "\x1b[?25h"
    cursorHome = "\x1b[H"
    cursorBlock = "\x1b[2 q"
    cursorUnderline = "\x1b[4 q"
    colorReset = "\x1b[m"
    colorInverted = "\x1b[7m"
    fontBold = "\x1b[1m"
)


type Renderer struct {
    // size
    nCols int
    nRows int

    // offset
    offsetX int
    offsetY int
    cursorX int
    cursorY int

    // output
    buffer *Buffer
}

func New(nRows, nCols int) (*Renderer, error) {
    return &Renderer{
        nCols: nCols,
        nRows: nRows - 2,
        buffer: &Buffer{},
    }, nil
}
func (r *Renderer) Restore() {
    r.buffer.Clear()
    r.buffer.Append(clearScreen)
    r.buffer.Append(cursorHome)
    r.buffer.Append(cursorShow)

    r.buffer.Write()
}
func (r *Renderer) Render(e *editor.Editor) {
    r.buffer.Clear()
    r.buffer.Append(cursorHide)
    r.buffer.Append(cursorHome)

    r.calculateOffset(e.CursorX, e.CursorY)

    r.drawRows(e.Content)
    r.drawStatusBar(e)
    r.drawCursor(e.Mode, len(e.Command))

    r.buffer.Write()
}

// =======
// private
// =======
func (r *Renderer) calculateOffset(editorCursorX, editorCursorY int) {
    renderStartY := r.offsetY
    renderEndY := r.offsetY + r.nRows

    if editorCursorY < renderStartY {
        r.offsetY = editorCursorY
    } else if editorCursorY >= renderEndY {
        r.offsetY = editorCursorY - r.nRows + 1
    }

    r.cursorX = editorCursorX
    r.cursorY = editorCursorY - r.offsetY
}
func (r *Renderer) drawRows(content []string) {
    contentLen := len(content)

    // edited rows
    for y := range contentLen {
        if y == r.nRows {
            break
        }

        row := content[y + r.offsetY]
        if len(row) > r.nCols {
            row = row[:r.nCols]
        }

        r.buffer.Append(row)
        r.buffer.Append(clearLine)
        r.buffer.Append("\r\n")
    }

    // unedited rows
    for y := contentLen; y < r.nRows; y++ {
        if y == r.nRows {
            break
        }
        
        r.buffer.Append(clearLine)
        r.buffer.Append("~\r\n")
    }
}
func (r *Renderer) drawStatusBar(e *editor.Editor) {
    r.buffer.Append(colorInverted)
    r.buffer.Append(fontBold)

    unsaved := ""
    if e.Dirty {
        unsaved = "[+]"
    }
    status := fmt.Sprintf(" %s | %s %s", e.Mode.String(), e.Filename, unsaved)
    for len(status) < r.nCols {
        status += " "
    }
    r.buffer.Append(status)

    r.buffer.Append(colorReset)
    r.buffer.Append(clearLine)
    r.buffer.Append("\r\n")

    // info bar
    r.buffer.Append(" " + e.Info)
    r.buffer.Append(clearLine)
}
func (r *Renderer) drawCursor(mode editor.EditorMode, commandLen int) {
    renderCursorX := r.cursorX + 1
    renderCursorY := r.cursorY + 1

    if mode == editor.Command {
        renderCursorX = commandLen + 2
        renderCursorY = r.nRows + 2
    }

    cursorPos := fmt.Sprintf("\x1b[%d;%dH", renderCursorY, renderCursorX)
    r.buffer.Append(cursorPos)

    if mode == editor.Normal { 
        r.buffer.Append(cursorBlock)
    } else {
        r.buffer.Append(cursorUnderline)
    }
    
    r.buffer.Append(cursorShow)
}
