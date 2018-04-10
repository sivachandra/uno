package lex

import (
	"bufio"
	"fmt"
	"os"
	"uno/lex/token"
)

func matchTokens(file string, ts token.KindSet, tokens []Token) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Error opening simple text file. \n%s", err.Error())
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var goEsr GoESR
	tz, err := NewTokenizer(br, ts, goEsr)
	if err != nil {
		return err
	}

	for _, exp := range tokens {
		if !tz.HasNext() {
			return fmt.Errorf("Expected a token at %s:%d:%d", file, exp.Line, exp.Col)
		}
		actual, err := tz.NextToken()
		if err != nil {
			return fmt.Errorf("Error at %d:%d: %s", tz.Line(), tz.Col(), err.Error())
		}

		if exp.Kind != actual.Kind {
			return fmt.Errorf(
				"Expected a token of kind %d at %s:%d:%d, got %d of value '%s'.",
				exp.Kind, file, exp.Line, exp.Col, actual.Kind, actual.Value)
		}

		if exp.Line != actual.Line {
			return fmt.Errorf("Expected line %d, but got %d for token with type %d.",
				exp.Line, actual.Line, actual.Kind)
		}

		if exp.Col != actual.Col {
			return fmt.Errorf("Expected col %d, but got %d.", exp.Col, actual.Col)
		}

		if string(actual.Value) != string(exp.Value) {
			return fmt.Errorf(
				"Expected token with value '%s' at %s:%d:%d, but got '%s'.",
				exp.Value, file, exp.Line, exp.Col, actual.Value)
		}
	}

	return nil
}
