package combinator

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/text/encoding"
)

type Combinator struct {
	buf *bytes.Buffer
	enc io.Writer
}

func New(e *encoding.Encoder) *Combinator {
	c := &Combinator{buf: new(bytes.Buffer)}
	if e == nil {
		c.enc = c.buf
	} else {
		c.enc = e.Writer(c.buf)
	}
	return c
}

func (c *Combinator) B(b ...byte) *Combinator {
	if _, err := c.buf.Write(b); err != nil {
		panic(err)
	}
	return c
}

func (c *Combinator) EB(b ...byte) *Combinator {
	if _, err := c.enc.Write(b); err != nil {
		panic(err)
	}
	return c
}

func (c *Combinator) S(s string) *Combinator {
	return c.B([]byte(s)...)
}

func (c *Combinator) ES(s string) *Combinator {
	return c.EB([]byte(s)...)
}

func (c *Combinator) FS(f string, args ...interface{}) *Combinator {
	return c.S(fmt.Sprintf(f, args...))
}

func (c *Combinator) EFS(f string, args ...interface{}) *Combinator {
	return c.ES(fmt.Sprintf(f, args...))
}

func (c *Combinator) Bytes() []byte {
	return c.buf.Bytes()
}

func (c *Combinator) String() string {
	return c.buf.String()
}
