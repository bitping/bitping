package parser

// QueryStatement is a structure that holds all query
// information based on blockchain events
type QueryStatement struct {
	Fields []string
	Values []string
}

// Parser represents a parser and it's current state
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buf size
	}
}
