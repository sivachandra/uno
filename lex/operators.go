package lex

import (
	"uno/lex/char"
	"uno/lex/token_kind"
)

var OperatorMap = map[string]uint32{
	"+":   token_kind.Add,
	"-":   token_kind.Sub,
	"*":   token_kind.Mul,
	"/":   token_kind.Div,
	"%":   token_kind.Mod,
	"=":   token_kind.Assign,
	"+=":  token_kind.AddAssign,
	"-=":  token_kind.SubAssign,
	"*=":  token_kind.MulAssign,
	"/=":  token_kind.DivAssign,
	"%=":  token_kind.ModAssign,
	"<<":  token_kind.LeftShift,
	">>":  token_kind.RightShift,
	"<<=": token_kind.LeftShiftAssign,
	">>=": token_kind.RightShiftAssign,
	"==":  token_kind.Equal,
	"===": token_kind.TripleEqual,
	"!=":  token_kind.NotEqual,
	"&":   token_kind.BitwiseAnd,
	"&=":  token_kind.BitwiseAndAssign,
	"|":   token_kind.BitwiseOr,
	"|=":  token_kind.BitwiseOrAssign,
	"^":   token_kind.BitwiseXor,
	"^=":  token_kind.BitwiseXorAssign,
	"~":   token_kind.BitwiseNeg,
	"~=":  token_kind.BitwiseNegAssign,
	"&&":  token_kind.LogicalAnd,
	"||":  token_kind.LogicalOr,
	"!":   token_kind.LogicalNot,
	".":   token_kind.Dot,
	"..":  token_kind.InclusiveRange,
	"...": token_kind.ExclusiveRange,
	",":   token_kind.Comma,
	"(":   token_kind.LeftParen,
	")":   token_kind.RightParen,
	"[":   token_kind.LeftBracket,
	"]":   token_kind.RightBracket,
	"{":   token_kind.LeftBrace,
	"}":   token_kind.RightBrace,
	":":   token_kind.Colon,
	"::":  token_kind.ScopeResolution,
	";":   token_kind.Semicolon,
	"<":   token_kind.LessThan,
	"<=":  token_kind.LessThanEqual,
	">":   token_kind.GreaterThan,
	">=":  token_kind.GreaterThanEqual,
	"->":  token_kind.ReturnArrow,
	"<-":  token_kind.ChannelIO,
	"++":  token_kind.UnaryIncrement,
	"--":  token_kind.UnaryDecrement,
	"**":  token_kind.MulPower,
}

func (tz *Tokenizer) hasCompAssign(op []rune) bool {
	cop := append(op, char.Equal)
	tt, e := OperatorMap[string(cop)]
	if !e {
		return false
	}

	return tz.ts.Contains(tt)
}

// Reads operators of different flavors, but all having the same starting
// character. Consider for example '<'. The argument |tt| refers to the
// plain '<' token, |tta| refers to '<=' token, 'ttr' refers to the '<<'
// token, and |ttra| refers to the '<<=' token_kind.
//
// tt - The basic token type of the single character.
// tta - The token type of the token which consists of the basic character
//     followed by the '=' character.
// ttr - The token type of the token which consists of the basic character
//    repeated twice.
// ttra - The token type of the token which consists of the |ttr| token
//    followed by the '=' character.
func (tz *Tokenizer) readOpFlavors(tt, tta, ttr, ttra uint32) (*Token, error) {
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
			token_kind.Add, token_kind.AddAssign, token_kind.UnaryIncrement, token_kind.Invalid)
	case char.Minus:
		// The operator '->' is not read by readOpFlavors. Hence we handle it
		// as a separate case.
		if tz.ts.Contains(token_kind.ReturnArrow) {
			cc, err := tz.r.PeekSlice(2)
			if err == nil && cc[1] == char.GreaterThan {
				s, err := tz.r.ReadSlice(2)
				if err != nil {
					return nil, err
				}
				return newToken(token_kind.ReturnArrow, s, line, col), nil
			}
		}
		return tz.readOpFlavors(
			token_kind.Sub, token_kind.SubAssign, token_kind.UnaryDecrement, token_kind.Invalid)
	case char.Mul:
		return tz.readOpFlavors(
			token_kind.Mul, token_kind.MulAssign, token_kind.MulPower, token_kind.Invalid)
	case char.Div:
		return tz.readOpFlavors(
			token_kind.Div, token_kind.DivAssign, token_kind.Invalid, token_kind.Invalid)
	case char.Mod:
		return tz.readOpFlavors(
			token_kind.Mod, token_kind.ModAssign, token_kind.Invalid, token_kind.Invalid)
	case char.Equal:
		return tz.readOpFlavors(
			token_kind.Assign, token_kind.Equal, token_kind.TripleEqual, token_kind.Invalid)
	case char.Tilde:
		return tz.readOpFlavors(
			token_kind.BitwiseNeg, token_kind.BitwiseNegAssign, token_kind.Invalid, token_kind.Invalid)
	case char.Carot:
		return tz.readOpFlavors(
			token_kind.BitwiseXor, token_kind.BitwiseXorAssign, token_kind.Invalid, token_kind.Invalid)
	case char.Exclaim:
		return tz.readOpFlavors(
			token_kind.LogicalNot, token_kind.NotEqual, token_kind.Invalid, token_kind.Invalid)
	case char.Ampersand:
		return tz.readOpFlavors(
			token_kind.BitwiseAnd, token_kind.BitwiseAndAssign, token_kind.LogicalAnd,
			token_kind.Invalid)
	case char.Pipe:
		return tz.readOpFlavors(
			token_kind.BitwiseOr, token_kind.BitwiseOrAssign, token_kind.LogicalOr,
			token_kind.Invalid)
	case char.Dot:
		return tz.readOpFlavors(
			token_kind.Dot, token_kind.Invalid, token_kind.InclusiveRange, token_kind.Invalid)
	case char.Comma:
		return tz.readOpFlavors(
			token_kind.Comma, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.LeftParen:
		return tz.readOpFlavors(
			token_kind.LeftParen, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.RightParen:
		return tz.readOpFlavors(
			token_kind.RightParen, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.LeftBracket:
		return tz.readOpFlavors(
			token_kind.LeftBracket, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.RightBracket:
		return tz.readOpFlavors(
			token_kind.RightBracket, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.LeftBrace:
		return tz.readOpFlavors(
			token_kind.LeftBrace, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.RightBrace:
		return tz.readOpFlavors(
			token_kind.RightBrace, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.Colon:
		return tz.readOpFlavors(
			token_kind.Colon, token_kind.Invalid, token_kind.ScopeResolution, token_kind.Invalid)
	case char.Semicolon:
		return tz.readOpFlavors(
			token_kind.Semicolon, token_kind.Invalid, token_kind.Invalid, token_kind.Invalid)
	case char.LessThan:
		return tz.readOpFlavors(
			token_kind.LessThan, token_kind.LessThanEqual, token_kind.LeftShift,
			token_kind.LeftShiftAssign)
	case char.GreaterThan:
		return tz.readOpFlavors(
			token_kind.GreaterThan, token_kind.GreaterThanEqual, token_kind.RightShift,
			token_kind.RightShiftAssign)
	}

	return nil, unExpectedCharacterError(c)
}
