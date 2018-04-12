package token_kind

const (
	Invalid = uint32(iota)

	// Keywords
	KeywordAnd
	KeywordAs
	KeywordAssert
	KeywordBreak
	KeywordClass
	KeywordConst
	KeywordContinue
	KeywordDef
	KeywordDel
	KeywordElif
	KeywordElse
	KeywordExcept
	KeywordCFalse
	KeywordPyFalse
	KeywordFinally
	KeywordFor
	KeywordFrom
	KeywordGlobal
	KeywordIf
	KeywordImport
	KeywordIn
	KeywordIs
	KeywordLambda
	KeywordNonlocal
	KeywordNot
	KeywordNull
	KeywordOr
	KeywordPass
	KeywordRaise
	KeywordReturn
	KeywordCTrue  // The boolean value 'true'
	KeywordPyTrue // The boolean value 'True'
	KeywordTry
	KeywordWhile
	KeywordWith
	KeywordYield

	// A C indentifier.
	Identifier

	// Python decorator:
	// @<IDENTIFIER>
	PythonDecorator

	// C preprocessor directive:
	// #<IDENTIFIER>
	CPPDirective

	DoubleQuoteString
	SingleQuoteString
	BackQuoteString
	PyMultilineString // """ ... """

	SingleQuoteCharacter // C style char

	DecimalInteger
	HexInteger
	OctInteger
	FloatNumber

	// C++ single line comment: //...
	CSingleLineComment
	// C++ multiline comment: /* ... */
	CMultiLineComment

	// Python single line comment: #...
	PySingleLineComment

	// The following tokens will also be used as operators.
	Add
	Sub
	Mul
	Div
	Mod

	BitwiseAnd
	BitwiseOr
	BitwiseXor
	BitwiseNot
	BitwiseNeg

	LeftShift
	RightShift

	Assign
	AddAssign
	SubAssign
	MulAssign
	DivAssign
	ModAssign
	LeftShiftAssign
	RightShiftAssign
	BitwiseAndAssign
	BitwiseOrAssign
	BitwiseXorAssign
	BitwiseNotAssign
	BitwiseNegAssign

	Equal
	TripleEqual // Ruby "===" operator
	NotEqual

	Dot
	InclusiveRange
	ExclusiveRange
	Comma

	LeftParen
	RightParen
	LeftBracket
	RightBracket
	LeftBrace
	RightBrace
	Colon
	ScopeResolution // Ruby and C++ scope resolution "::"
	Semicolon

	LogicalAnd
	LogicalOr
	LogicalNot

	LessThan
	LessThanEqual
	GreaterThan
	GreaterThanEqual

	ReturnArrow
	ChannelIO

	UnaryIncrement
	UnaryDecrement
	MulPower

	// Indent at the beginning of a line
	Indent

	// In some languages new line and tab are relevant
	NewLine
	Tab

	LineJoin

	FirstInvalidTokenKind
)

