package lex

import (
	"fmt"
	"uno/lex/char"
	"uno/lex/token"
)

func isDecimalDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isOctDigit(c rune) bool {
	return '0' <= c && c <= '7'
}

func isHexDigit(c rune) bool {
	return isDecimalDigit(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func decimalSyntaxError() error {
	return fmt.Errorf("Bad decimal integer syntax.")
}

func hexSyntaxError() error {
	return fmt.Errorf("Bad hex integer syntax.")
}

func floatSyntaxError() error {
	return fmt.Errorf("Bad floating point number syntax error.")
}

func octSyntaxError() error {
	return fmt.Errorf("Bad octal integer syntax.")
}

func (tz *Tokenizer) readCharAfterE() (rune, error) {
	c, err := tz.r.ReadChar()
	if err != nil {
		return 0, err
	}
	if !(isDecimalDigit(c) || c == '+' || c == '-') {
		return 0, floatSyntaxError()
	}
	return c, nil
}

func (tz *Tokenizer) readNumber() (*Token, error) {
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	c, err := tz.r.ReadChar()
	if err != nil {
		return nil, err
	}
	if !isDecimalDigit(c) && c != char.Dot {
		// This function should be called only when a number is expected.
		return nil, fmt.Errorf("Unexpected '%c' while reading a number.", c)
	}

	n := []rune{c}
	dec := false
	hex := false
	oct := false
	float := false
	// Flag to check whether an 'e' or 'E' of a float number has been read.
	exp := false

	if c == '0' {
		// The number is either octal, hex, float, or just plain
		// decimal zero
		c, err = tz.r.PeekChar()
		if err != nil {
			// Decimal zero at the end of the stream.
			return tz.newValidToken(token.DecimalInteger, n, line, col)
		}

		if isAnyWhiteSpace(c) || (IsDelimiter(c) && c != char.Dot) {
			// Decimal zero
			return tz.newValidToken(token.DecimalInteger, n, line, col)
		}

		c, err := tz.r.ReadChar()
		if err != nil {
			return nil, err
		}
		// Storing it before error detection facilitates
		// reading subsquent characters.
		n = append(n, c)

		switch {
		case c == 'x' || c == 'X':
			if !tz.ts.Contains(token.HexInteger) {
				return nil, unExpectedCharacterError(c)
			}
			hex = true
		case c == char.Dot:
			float = true
		case c == 'E' || c == 'e':
			float = true
			// If it is an 'e' or 'E', then it has to be followed by a
			// '+', '-' or a decimal digit.
			c, err = tz.readCharAfterE()
			if err != nil {
				return nil, err
			}
			n = append(n, c)
			exp = true
		case isDecimalDigit(c):
			// It is either an octal number or a float number.
			// We will treat it as octal until we find a '.' or 'e' or 'E'.
			oct = true
		default:
			return nil, unExpectedCharacterError(c)
		}
	} else if c == char.Dot {
		float = true
	} else {
		// Start with decimal. Change to float '.', 'e' or 'E' is
		// encountered.
		dec = true
	}

	// Since a number starting with '0' can also be a float or octal
	// number, we will treat it as octal until we find a '.', 'e' or 'E'.
	// However, we will keep track of invalid octal digits while doing so.
	invalidOct := false

	for {
		c, err = tz.r.PeekChar()
		if err != nil || isAnyWhiteSpace(c) || (IsDelimiter(c) && c != char.Dot) {
			break
		}

		c, err = tz.r.ReadChar()
		if err != nil {
			return nil, err
		}
		// Storing it before error detection facilitates
		// reading subsquent characters.
		n = append(n, c)

		switch {
		case hex:
			if !isHexDigit(c) {
				return nil, hexSyntaxError()
			}
		case oct:
			if c == char.Dot || c == 'E' || c == 'e' {
				oct = false
				float = true
				exp = !(c == char.Dot)
			}

			if !isDecimalDigit(c) {
				return nil, octSyntaxError()
			}
			if c == '8' || c == '9' {
				invalidOct = true
			}
		case dec:
			switch {
			case c == char.Dot:
				dec = false
				float = true
			case c == 'E' || c == 'e':
				// If it is an 'e' or 'E', then it has to be followed by a
				// '+', '-' or a decimal digit.
				c, err = tz.readCharAfterE()
				if err != nil {
					return nil, err
				}
				dec = false
				float = true
				exp = true
				n = append(n, c)
			case !isDecimalDigit(c):
				return nil, decimalSyntaxError()
			}
		case float:
			// This case is either hit after reading a '.', 'e', 'E'.
			if c == 'E' || c == 'e' {
				if exp {
					return nil, floatSyntaxError()
				}
				// If it is an 'e' or 'E', then it has to be followed by a
				// '+', '-' or a decimal digit.
				c, err = tz.readCharAfterE()
				if err != nil {
					return nil, err
				}
				exp = true
				n = append(n, c)
			} else if !isDecimalDigit(c) {
				return nil, floatSyntaxError()
			}
		default:
			// Read the character as part of the number.
			break
		}
	}

	var tt token.Kind
	if dec {
		tt = token.DecimalInteger
	}
	if hex {
		if len(n) < 3 {
			return nil, hexSyntaxError()
		}
		tt = token.HexInteger
	}
	if float {
		if len(n) < 2 {
			return nil, floatSyntaxError()
		}
		tt = token.FloatNumber
	}
	if oct {
		if len(n) < 2 {
			return nil, octSyntaxError()
		}
		// The oct integer string was constructed even if
		// the digits were 8 or 9 for the sake of
		// convenience in the above flow. Check for validity
		// now.
		if invalidOct {
			return nil, octSyntaxError()
		}
		tt = token.OctInteger
	}

	return newToken(tt, n, line, col), nil
}
