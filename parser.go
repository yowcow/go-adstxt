package adstxt

import (
	"bufio"
	"io"
	"strings"
)

// Parser is an iterative ads.txt parser
type Parser struct {
	scanner *bufio.Scanner
}

// NewParser returns a Parser
func NewParser(r io.Reader) *Parser {
	return &Parser{bufio.NewScanner(r)}
}

// Parse returns a *Record or an error
// Deprecated: use ParseRow instead.
func (p *Parser) Parse() (*Record, error) {
	// scans for the first valid row or returns an error otherwise
	for p.scanner.Scan() {
		text := normalizeRow(p.scanner.Text())

		// blank line
		if len(text) == 0 {
			continue
		}

		// returns when either is non-nil
		r, err := parseRow(text)
		if err != nil {
			return nil, err
		}
		if r != nil && r.Record != nil {
			return r.Record, nil
		}
	}

	if err := p.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}

// ParseRow returns a Row
func (p *Parser) ParseRow() (*Row, error) {
	// scans for the first valid row or returns an error otherwise
	for p.scanner.Scan() {
		text := normalizeRow(p.scanner.Text())

		// blank line
		if len(text) == 0 {
			continue
		}

		r, err := parseRow(text)
		if r != nil || err != nil {
			return r, err
		}

	}

	if err := p.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}

func normalizeRow(s string) string {
	text := strings.TrimSpace(s)

	// remove comment
	if idx := strings.IndexRune(text, '#'); idx >= 0 {
		text = text[0:idx]
	}

	return text
}
