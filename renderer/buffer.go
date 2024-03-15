package renderer

import (
    "os"
)

type Buffer struct {
    content string
}
func (b *Buffer) Append(s string) {
    b.content += s
}
func (b *Buffer) Clear() {
    b.content = ""
}
func (b *Buffer) String() string {
    return b.content
}
func (b *Buffer) Write() {
    os.Stdout.WriteString(b.content)
}


