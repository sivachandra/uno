package lex

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestSimpleTextRead(t *testing.T) {
	file, err := os.Open("test_data/simple_text")
	if err != nil {
		t.Errorf("Error opening simple text file. \n%s", err.Error())
	}
	defer file.Close()

	br := bufio.NewReader(file)
	r := NewCharReader(br)

	if r.Line() != 0 {
		t.Errorf("The line number before reading any char is not 0.")
	}
	if r.Col() != 0 {
		t.Errorf("The column number before reading any char is not 0.")
	}
	if r.NextLine() != 1 {
		t.Errorf("Next line before reading any char is not 1.")
	}
	if r.NextCol() != 1 {
		t.Errorf("Next column before reading any char is not 1.")
	}
	if !r.PreviousWasNewLine() {
		t.Errorf("The last char read is not a new line before reading any char.")
	}

	c, err := r.ReadChar()
	if err != nil {
		t.Errorf("Error reading char at %d:%d.", r.Line(), r.Col())
	}
	if c != 'h' {
		t.Errorf("Expected 'h' at %d:%d.", r.Line(), r.Col())
	}
	if r.Line() != 1 {
		t.Errorf("Line number of the last read char is not 1.")
	}
	if r.Col() != 1 {
		t.Errorf("Column number of the last read char is not 1.")
	}

	s, err := r.ReadSlice(11)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(s) != "ello, world" {
		t.Errorf("Expected 'ello, world', but got '%s'.", string(s))
	}
	if r.Line() != 1 {
		t.Errorf("Line after reading out 'hello, world' is not 1.")
	}
	if r.Col() != 12 {
		t.Errorf("Column after reading out 'hello, world' is not 12.")
	}

	c, err = r.ReadChar()
	if err != nil {
		t.Errorf(err.Error())
	}
	if c != '\n' {
		t.Errorf("Expected new line but got '%c' at %d:%d.", c, r.Line(), r.Col())
	}
	if !r.PreviousWasNewLine() {
		t.Errorf("PreviousWasNewLine() does not return true after reading out a new line.")
	}
	if r.Line() != 1 {
		t.Errorf("Expected Line() to return 1 after reading the first new line.")
	}
	if r.Col() != 13 {
		t.Errorf("Expected the column of the first new line to be 13.")
	}
	if r.NextLine() != 2 {
		t.Errorf("NextLine() does not return 2 after reading the first line.")
	}

	c, err = r.ReadChar()
	if err != nil {
		t.Errorf(err.Error())
	}
	if c != '\n' {
		t.Errorf("Expected new line but got '%c' at %d:%d.", c, r.Line(), r.Col())
	}
	if !r.PreviousWasNewLine() {
		t.Errorf("PreviousWasNewLine() does not return true after reading out a new line.")
	}
	if r.Line() != 2 {
		t.Errorf("Expected Line() to return 2 after reading the second new line.")
	}
	if r.Col() != 1 {
		t.Errorf("Expected the column of the second new line to be 1.")
	}
	if r.NextLine() != 3 {
		t.Errorf("NextLine() does not return 3 after reading the second new line.")
	}

	c, err = r.PeekChar()
	if err != nil {
		t.Errorf(err.Error())
	}
	if c != 'I' {
		t.Errorf("Expected 'I' but got '%c' at %d:%d.", c, r.Line(), r.Col())
	}

	s, err = r.PeekSlice(11)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(s) != "I am a file" {
		t.Errorf("Expected 'I am a file', but got '%s'.", string(s))
	}

	s, err = r.ReadSlice(43)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(s) != "I am a file you can use to test CharReader." {
		t.Errorf(
			"Expected 'I am file you can use to test CharReader.', but got '%s'.",
			string(s))
	}
	if r.Line() != 3 {
		t.Errorf("Expected Line() to return 3, but returns %d.", r.Line())
	}
	if r.Col() != 43 {
		t.Errorf("Expected Col() to return 43, but returns %d.", r.Col())
	}

	c, err = r.ReadChar()
	if err != nil {
		t.Errorf(err.Error())
	}
	if c != '\n' {
		t.Errorf("Expected new line but got '%c' at %d:%d.", c, r.Line(), r.Col())
	}
	if !r.PreviousWasNewLine() {
		t.Errorf("PreviousWasNewLine() does not return true after reading out a new line.")
	}
	if r.Line() != 3 {
		t.Errorf("Expected Line() to return 3 after reading the third new line.")
	}
	if r.Col() != 44 {
		t.Errorf("Expected the column of the third new line to be 44.")
	}
	if r.NextLine() != 4 {
		t.Errorf("NextLine() does not return 4 after reading the third new line.")
	}

	c, err = r.ReadChar()
	if err != nil {
		t.Errorf(err.Error())
	}
	c, err = r.ReadChar()
	if err != io.EOF {
		t.Errorf("Expected EOF when attempting to read beyong the last new line.")
	}
}

func TestReadingBeyongEOF(t *testing.T) {
	file, err := os.Open("test_data/simple_text")
	if err != nil {
		t.Errorf("Error opening simple text file. \n%s", err.Error())
	}
	defer file.Close()

	br := bufio.NewReader(file)
	r := NewCharReader(br)

	_, err = r.PeekSlice(1000)
	if err != io.EOF {
		t.Errorf("Peeking too much should result in EOF error.")
	}

	s, err := r.ReadSlice(12)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(s) != "hello, world" {
		t.Errorf("Expected to read 'hello, world', but got '%s'.", string(s))
	}

	_, err = r.ReadSlice(1000)
	if err != io.EOF {
		t.Errorf("Reading too much should result in EOF error.")
	}
}
