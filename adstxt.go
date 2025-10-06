// Package adstxt implements Ads.txt protocol defined by iab.
package adstxt

import (
	"io"
	"net/http"
)

// Deprecated: to be removed
func Get(rawurl string) ([]Record, error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return Parse(resp.Body)
}

// Deprecated: use ParseRows instead.
func Parse(in io.Reader) ([]Record, error) {
	records := make([]Record, 0)
	p := NewParser(in)

LOOP:
	for {
		r, err := p.Parse()
		if err == io.EOF {
			break LOOP
		}
		if err != nil {
			return nil, err
		}
		if r != nil {
			records = append(records, *r)
		}
	}

	return records, nil
}

func ParseRows(in io.Reader) ([]Row, error) {
	rows := make([]Row, 0)
	p := NewParser(in)

	for {
		row, err := p.ParseRow()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if row != nil {
			rows = append(rows, *row)
		}
	}

	return rows, nil
}
