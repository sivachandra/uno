package lex

import (
	"uno/lex/char"
	"uno/lex/token"
)

var Operators = map[string]token.Kind{
	"+":   token.Add,
	"-":   token.Sub,
	"*":   token.Mul,
	"/":   token.Div,
	"%":   token.Mod,
	"=":   token.Assign,
	"+=":  token.AddAssign,
	"-=":  token.SubAssign,
	"*=":  token.MulAssign,
	"/=":  token.DivAssign,
	"%=":  token.ModAssign,
	"<<":  token.LeftShift,
	">>":  token.RightShift,
	"<<=": token.LeftShiftAssign,
	">>=": token.RightShiftAssign,
	"==":  token.Equal,
	"===": token.TripleEqual,
	"!=":  token.NotEqual,
	"&":   token.BitwiseAnd,
	"&=":  token.BitwiseAndAssign,
	"|":   token.BitwiseOr,
	"|=":  token.BitwiseOrAssign,
	"^":   token.BitwiseXor,
	"^=":  token.BitwiseXorAssign,
	"~":   token.BitwiseNeg,
	"~=":  token.BitwiseNegAssign,
	"&&":  token.LogicalAnd,
	"||":  token.LogicalOr,
	"!":   token.LogicalNot,
	".":   token.Dot,
	"..":  token.InclusiveRange,
	"...": token.ExclusiveRange,
	",":   token.Comma,
	"(":   token.LeftParen,
	")":   token.RightParen,
	"[":   token.LeftBracket,
	"]":   token.RightBracket,
	"{":   token.LeftBrace,
	"}":   token.RightBrace,
	":":   token.Colon,
	"::":  token.ScopeResolution,
	";":   token.Semicolon,
	"<":   token.LessThan,
	"<=":  token.LessThanEqual,
	">":   token.GreaterThan,
	">=":  token.GreaterThanEqual,
	"->":  token.ReturnArrow,
	"<-":  token.ChannelIO,
	"++":  token.UnaryIncrement,
	"--":  token.UnaryDecrement,
	"**":  token.MulPower,
}

func (tz *Tokenizer) hasCompAssign(op []rune) bool {
	cop := append(op, char.Equal)
	tt, e := Operators[string(cop)]
	if !e {
		return false
	}

	return tz.ts.Contains(tt)
}

// Reads operators of different flavors, but all having the same starting
// character. Consider for example '<'. The argument |tt| refers to the
// plain '<' token, |tta| refers to '<=' token, 'ttr' refers to the '<<'
// token, and |ttra| refers to the '<<=' token.
//
// tt - The basic token type of the single character.
// tta - The token type of the token which consists of the basic character
//     followed by the '=' character.
// ttr - The token type of the token which consists of the basic character
//    repeated twice.
// ttra - The token type of the token which consists of the |ttr| token
//    followed by the '=' character.
func (tz *Tokenizer) readOpFlavors(tt, tta, ttr, ttra token.Kind) (*Token, error) {
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	c, err := tz.r.ReadChar()
	if err != nil {
		return nil, err
	}

	op := []rune{c}
	c1, err := tz.r.PeekChar()
	switch {
	case err != nil:
		return tz.newValidToken(tt, op, line, col)
	case c1 == char.Equal && tz.hasCompAssign(op):
		c1, err = tz.r.ReadChar()
		if err != nil {
			return nil, err
		}
		op = append(op, c1)
		return newToken(tta, op, line, col), nil
	case c1 == c && tz.ts.Contains(ttr):
		c1, err = tz.r.ReadChar()
		if err != nil {
			return nil, err
		}
		op = append(op, c1)

		c2, err := tz.r.PeekChar()
		if err != nil || c2 != char.Equal || !tz.hasCompAssign(op) {
			return newToken(ttr, op, line, col), nil
		}
		c2, err = tz.r.ReadChar()
		if err != nil {
			return nil, err
		}
		op = append(op, c2)
		return newToken(ttra, op, line, col), nil
	default:
		return tz.newValidToken(tt, op, line, col)
	}
}

func (tz *Tokenizer) readOperator() (*Token, error) {
	c, err := tz.r.PeekChar()
	if err != nil {
		return nil, err
	}

	line := tz.r.NextLine()
	col := tz.r.NextCol()

	switch c {
	case char.Plus:
		return tz.readOpFlavors(
			token.Add, token.AddAssign, token.UnaryIncrement, token.Invalid)
	case char.Minus:
		// The operator '->' is not read by readOpFlavors. Hence we handle it
		// as a separate case.
		if tz.ts.Contains(token.ReturnArrow) {
			cc, err := tz.r.PeekSlice(2)
			if err == nil && cc[1] == char.GreaterThan {
				s, err := tz.r.ReadSlice(2)
				if err != nil {
					return nil, err
				}
				return newToken(token.ReturnArrow, s, line, col), nil
			}
		}
		return tz.readOpFlavors(
			token.Sub, token.SubAssign, token.UnaryDecrement, token.Invalid)
	case char.Mul:
		return tz.readOpFlavors(
			token.Mul, token.MulAssign, token.MulPower, token.Invalid)
	case char.Div:
		return tz.readOpFlavors(
			token.Div, token.DivAssign, token.Invalid, token.Invalid)
	case char.Mod:
		return tz.readOpFlavors(
			token.Mod, token.ModAssign, token.Invalid, token.Invalid)
	case char.Equal:
		return tz.readOpFlavors(
			token.Assign, token.Equal, token.TripleEqual, token.Invalid)
	case char.Tilde:
		return tz.readOpFlavors(
			token.BitwiseNeg, token.BitwiseNegAssign, token.Invalid, token.Invalid)
	case char.Carot:
		return tz.readOpFlavors(
			token.BitwiseXor, token.BitwiseXorAssign, token.Invalid, token.Invalid)
	case char.Exclaim:
		return tz.readOpFlavors(
			token.LogicalNot, token.NotEqual, token.Invalid, token.Invalid)
	case char.Ampersand:
		return tz.readOpFlavors(
			token.BitwiseAnd, token.BitwiseAndAssign, token.LogicalAnd,
			token.Invalid)
	case char.Pipe:
		return tz.readOpFlavors(
			token.BitwiseOr, token.BitwiseOrAssign, token.LogicalOr,
			token.Invalid)
	case char.Dot:
		return tz.readOpFlavors(
			token.Dot, token.Invalid, token.InclusiveRange, token.Invalid)
	case char.Comma:
		return tz.readOpFlavors(
			token.Comma, token.Invalid, token.Invalid, token.Invalid)
	case char.LeftParen:
		return tz.readOpFlavors(
			token.LeftParen, token.Invalid, token.Invalid, token.Invalid)
	case char.RightParen:
		return tz.readOpFlavors(
			token.RightParen, token.Invalid, token.Invalid, token.Invalid)
	case char.LeftBracket:
		return tz.readOpFlavors(
			token.LeftBracket, token.Invalid, token.Invalid, token.Invalid)
	case char.RightBracket:
		return tz.readOpFlavors(
			token.RightBracket, token.Invalid, token.Invalid, token.Invalid)
	case char.LeftBrace:
		return tz.readOpFlavors(
			token.LeftBrace, token.Invalid, token.Invalid, token.Invalid)
	case char.RightBrace:
		return tz.readOpFlavors(
			token.RightBrace, token.Invalid, token.Invalid, token.Invalid)
	case char.Colon:
		return tz.readOpFlavors(
			token.Colon, token.Invalid, token.ScopeResolution, token.Invalid)
	case char.Semicolon:
		return tz.readOpFlavors(
			token.Semicolon, token.Invalid, token.Invalid, token.Invalid)
	case char.LessThan:
		return tz.readOpFlavors(
			token.LessThan, token.LessThanEqual, token.LeftShift,
			token.LeftShiftAssign)
	case char.GreaterThan:
		return tz.readOpFlavors(
			token.GreaterThan, token.GreaterThanEqual, token.RightShift,
			token.RightShiftAssign)
	}

	return nil, unExpectedCharacterError(c)
}
