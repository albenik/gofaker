package fake

import (
	"github.com/albenik/gofaker/fakereader"
	"github.com/albenik/gofaker/fakewriter"
)

type ReadWriter struct {
	fakereader.Reader
	fakewriter.Writer
}
