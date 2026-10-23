package writer

import (
	"bufio"
	"io"

	"github.com/Luckny/recl/internal/parser"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: bufio.NewWriter(w)}
}

func (w *Writer) Write(v parser.Value) {
	resp := parser.RESPValue(v)
	w.writer.Write([]byte(resp))
}

func (w *Writer) Flush() {
	w.writer.(*bufio.Writer).Flush()
}
