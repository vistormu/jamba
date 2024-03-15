package inputreader

import (
    "os"
    "jamba/key"
)

type InputReader struct {
    buffer []byte
}

func New() *InputReader {
    return &InputReader{make([]byte, 1)}
}

func (ir *InputReader) ReadKey() (key.Key, error) {
    for {
        readLen, _ := os.Stdin.Read(ir.buffer)
        if readLen == 1 {
            break
        }
    }

    if ir.buffer[0] == 27 {
        seq := make([]byte, 2)
        os.Stdin.Read(seq)
        key, err := key.New([3]byte{ir.buffer[0], seq[0], seq[1]})
        if err != nil {
            return 0, err
        }
        return key, nil
    }
    
    return key.New(ir.buffer[0])
}

