package main

import (
    "io"
    "os"
    "fmt"
    "jamba/terminal"
    "jamba/inputreader"
    "jamba/editor"
    "jamba/renderer"
)

func restore(t *terminal.Terminal, r *renderer.Renderer) {
    t.Restore()
    r.Restore()
}

func exit(err error, t *terminal.Terminal, r *renderer.Renderer) {
    restore(t, r)
    if err == io.EOF {
        os.Exit(0)
    }
    panic(err)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: jamba <filename>")
        os.Exit(1)
    }

    // terminal
    t, err := terminal.New()
    if err != nil {
        exit(err, t, nil)
    }

    // input reader
    ir := inputreader.New()

    // editor
    e, err := editor.New(os.Args[1])
    if err != nil {
        exit(err, t, nil)
    }

    // renderer
    r, err := renderer.New(t.NRows, t.NCols)
    if err != nil {
        exit(err, t, r)
    }

    defer restore(t, r)

    for {
        // render editor content
        r.Render(e)

        // get key
        key, err := ir.ReadKey()
        if err != nil {
            exit(err, t, r)
        }

        // update editor content
        err = e.Update(key)
        if err != nil {
            exit(err, t, r)
        }
    }
}
