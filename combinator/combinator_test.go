package combinator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"

	"github.com/albenik/gofaker/combinator"
)

var win1251 = charmap.Windows1251.NewEncoder()

func TestCombinator_B(t *testing.T) {
	c := combinator.New(nil).
		B(0x01, 0x02, 0x03).
		B('1', '2', '3')
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x31, 0x32, 0x33}, c.Bytes())
	assert.Equal(t, "\x01\x02\x03123", c.String())
}

func TestCombinator_EB_Windows1251(t *testing.T) {
	c := combinator.New(win1251).
		B(0x01, 0x02, 0x03, 0x04).
		EB(0x01, 0x02, 0x03, 0x04).
		B(0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4).
		EB(0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4).
		EB([]byte("ЧЫФ")...)

	assert.Equal(t, []byte{
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
		0xD7, 0xDB, 0xD4,
		0xD7, 0xDB, 0xD4,
	}, c.Bytes())
}

func TestCombinator_EB_None(t *testing.T) {
	c := combinator.New(nil).
		B(0x01, 0x02, 0x03, 0x04).
		EB(0x01, 0x02, 0x03, 0x04).
		B(0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4).
		EB(0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4).
		EB([]byte("ЧЫФ")...)

	assert.Equal(t, []byte{
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
	}, c.Bytes())
}

func TestCombinator_S(t *testing.T) {
	c := combinator.New(win1251).
		S("1234ЧЫФ").
		FS("%s%d", "ЧЫФ", 5678)

	assert.Equal(t, []byte{
		0x31, 0x32, 0x33, 0x34,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
		0x35, 0x36, 0x37, 0x38,
	}, c.Bytes())

	assert.Equal(t,
		"1234\xD0\xA7\xD0\xAB\xD0\xA4\xD0\xA7\xD0\xAB\xD0\xA45678",
		c.String(),
	)
}

func TestCombinator_ES_Windows1251(t *testing.T) {
	c := combinator.New(win1251).
		ES("1234ЧЫФ").
		EFS("%s%d", "ЧЫФ", 5678)

	assert.Equal(t, []byte{
		0x31, 0x32, 0x33, 0x34,
		0xD7, 0xDB, 0xD4,
		0xD7, 0xDB, 0xD4,
		0x35, 0x36, 0x37, 0x38,
	}, c.Bytes())

	assert.Equal(t,
		"1234\xD7\xDB\xD4\xD7\xDB\xD45678",
		c.String(),
	)
}

func TestCombinator_ES_None(t *testing.T) {
	c := combinator.New(nil).
		ES("1234ЧЫФ").
		EFS("%s%d", "ЧЫФ", 5678)

	assert.Equal(t, []byte{
		0x31, 0x32, 0x33, 0x34,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
		0xD0, 0xA7, 0xD0, 0xAB, 0xD0, 0xA4,
		0x35, 0x36, 0x37, 0x38,
	}, c.Bytes())

	assert.Equal(t,
		"1234\xD0\xA7\xD0\xAB\xD0\xA4\xD0\xA7\xD0\xAB\xD0\xA45678",
		c.String(),
	)
}
