package transport

import (
	"bufio"
	"io"
)

type rawHandler struct {
	Handler

	// in is the file descriptor that we should read raw commands from.
	in *bufio.Reader

	// out is the file descriptor that we should write raw output to.
	out *bufio.Writer
}

// NewRawHandler produces a transport handler that communicates over a Reader and a Writer,
// with the minimum amount of formatting.  This could be used for telnet-based connections
// or even stdin/stdout.
func NewRawHandler(in io.Reader, out io.Writer) Handler {
	l := rawHandler{
		in:  bufio.NewReader(in),
		out: bufio.NewWriter(out),
	}
	return &l
}

// ParseLine reads a message from the `in` Reader.
func (l *rawHandler) ReadLine() (string, error) {
	return l.in.ReadString('\n')
}

// WriteMsg writes a raw message out to the `out` Reader.
func (l *rawHandler) WriteMsg(b []byte) error {
	_, err := l.out.Write(b)
	return err
}
