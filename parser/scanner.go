package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the buffered reader
// Returns the rune(0) if an error ocurrs or io.EOF is returned)
func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	// If we see a whitelist then consume all the next whitespace
	// If we see a letter consume as an ident or reserved word
	// if we see a digit consume as a number
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isAddress(s, ch) {
		return s.scanAddress()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) {
		s.unread()
		return s.scanDigit()
	}

	// otherwise read individual character
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTRISK, string(ch)
	case ',':
		return COMMA, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequence whitelist char into the buffer
	// non-whitepsace characters and EOF exits the loop
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes
func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every next ident char
	// non-ident characters and EOF exits loop
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	}

	return IDENT, buf.String()
}

// scanAddress consumes the entire address and all contiguous address rules
func (s *Scanner) scanAddress() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteString("0")

	// For every next ident char
	// except non-ident chars and EOF exits loop
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return ADDRESS, buf.String()
}

// scanDigit consumes a digit, including periods until either
// rules don't apply or eof is consumed
func (s *Scanner) scanDigit() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// For all values, except eof and non-digits
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isDigit(ch) && !isExponential(ch) {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return DIGITS, buf.String()
}

// read reads the next rune from the buffered reader
// returns the rune(0)  if an error occus or EOF is returned
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on reader
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// isWhitelist returns true if run is a space, tab, newline
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isAddress(s *Scanner, a rune) bool {
	if a != '0' {
		return false
	} else if b := s.read(); b != eof {
		s.unread()
		return (a == '0' && b == 'x')
	}
	return false
}

// isLetter returns true if rune is a letter
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z')
}

// isDigit returns true if the rune is a digit
func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isExponential(ch rune) bool {
	return ch == 'e'
}

func isPeriod(ch rune) bool {
	return ch == '.'
}

var eof = rune(0)
