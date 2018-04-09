package lex

import (
	"fmt"
	"io"
	"unicode"
	"uno/lex/char"
)

type CharReader struct {
	cache   []rune
	r       io.RuneReader
	line    uint32
	col     uint32
	newLine bool
}

func NewCharReader(r io.RuneReader) *CharReader {
	cr := new(CharReader)
	cr.r = r
	cr.line = 0
	cr.col = 0
	cr.newLine = true

	return cr
}

// Returns the line on which the last successfully read or attempted
// character was present on. A value of 0 is returned before the first
// character is read.
func (r *CharReader) Line() uint32 {
	return r.line
}

// Returns the line on which the next character to be read exists.
func (r *CharReader) NextLine() uint32 {
	if r.newLine {
		return r.line + 1
	} else {
		return r.line
	}
}

// Returns the column on which the last successfully read or attempted
// character was present at. A value of 0 is returned before the first
// character is read.
func (r *CharReader) Col() uint32 {
	return r.col
}

// Returns the column on which the next character to be read exists.
func (r *CharReader) NextCol() uint32 {
	if r.newLine {
		return 1
	} else {
		return r.col + 1
	}
}

// Returns true if the last character read was a new line character.
func (r *CharReader) PreviousWasNewLine() bool {
	return r.newLine
}

func (r *CharReader) readOutChar() (rune, error) {
	c, s, err := r.r.ReadRune()
	if err != nil {
		return 0, err
	}
	if c == unicode.ReplacementChar && s == 1 {
		return 0, fmt.Errorf("Invalid unicode character.")
	}
	return c, nil
}

// Read a character and return it.
// If an error occurs while reading, it is not guaranteed to be
// recoverable. In general, it is a good idea to call the method
// PeekChar or PeekSlice before calling this method. There could
// be situations while reading a token wherein another character
// is required for the token to be valid. In such a situation,
// one can call ReadChar but should be aware the error is not
// guaranteed to be recoverable.
func (r *CharReader) ReadChar() (rune, error) {
	if r.newLine {
		r.col = 1
		r.line += 1
	} else {
		r.col += 1
	}

	var c rune
	if len(r.cache) > 0 {
		c = r.cache[0]
		r.cache = r.cache[1:]
	} else {
		var err error
		c, err = r.readOutChar()
		if err != nil {
			return 0, err
		}
	}

	if c == char.NewLine {
		r.newLine = true
	} else {
		r.newLine = false
	}

	return c, nil
}

// Reads and returns a slice on n characters.
// Recoverability from an error is not guaranteed. See the
// description of ReadChar for more information.
func (r *CharReader) ReadSlice(n uint32) ([]rune, error) {
	var s []rune
	for i := uint32(0); i < n; i++ {
		c, err := r.ReadChar()
		if err != nil {
			return nil, err
		}

		s = append(s, c)
	}

	return s, nil
}

// Peek and return a character if available.
// The error returned will be io.EOF if trying to peek
// beyond the end of the data. A rune value of 0 is
// returned on error.
func (r *CharReader) PeekChar() (rune, error) {
	if len(r.cache) > 0 {
		return r.cache[0], nil
	}

	c, err := r.readOutChar()
	if err != nil {
		return c, err
	}
	r.cache = append(r.cache, c)

	return c, nil
}

// Peek and return a slice of n characters if available.
// The error returned will be io.EOF if trying to peek
// beyond the end of the data. 'nil' is returned on error.
func (r *CharReader) PeekSlice(n uint32) ([]rune, error) {
	m := uint32(len(r.cache))
	for i := m; i < n; i++ {
		c, err := r.readOutChar()
		if err != nil {
			return nil, err
		}
		r.cache = append(r.cache, c)
	}
	return r.cache[0:n], nil
}
