package parser

// Token is a lexable token
type Token int

const (
	// ILLEGAL happens only if the parsable input is invalid
	ILLEGAL Token = iota
	// ADDRESS is a watchable address
	ADDRESS
	// VALUE is a watchable value
	VALUE
	// DIGITS is an actual number value
	DIGITS

	// EOF is end of file
	EOF
	// NL is a newline
	NL
	// WS is whitespace
	WS

	// IDENT an identifier
	IDENT

	// ASTRISK represents anything
	ASTRISK
	// COMMA represents and _and_ statement
	COMMA

	// SELECT keys to pass into the event
	SELECT
	// AND groups multiple conditional statements
	AND
	// EQUAL looks if items are equal
	EQUAL
	// NOT says an item is not something
	NOT
	// LESSTHAN says item is less than
	LESSTHAN
	// GREATERTHAN says an item is greater than
	GREATERTHAN
	// ORDER sets the ordering
	ORDER
)
