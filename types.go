package adstxt

import (
	"fmt"
	"regexp"
	"strings"
)

type Row interface {
	fmt.Stringer
}

var (
	_ Row = (*Record)(nil)
	_ Row = (*Variable)(nil)
)

const (
	AccountDirect AccountType = iota
	AccountReseller
	AccountOther
)

type AccountType int

// Record is ads.txt data field defined in iab.
type Record struct {
	// ExchangeDomain is domain name of the advertising system
	ExchangeDomain string

	// PublisherAccountID is the identifier associated with the seller
	// or reseller account within the advertising system.
	PublisherAccountID string

	// AccountType is an enumeration of the type of account.
	AccountType AccountType

	// AuthorityID is an ID that uniquely identifies the advertising system
	// within a certification authority.
	AuthorityID string
}

var (
	leadingBlankRe  = regexp.MustCompile(`\A[\s\t]+`)
	trailingBlankRe = regexp.MustCompile(`[\s\t]+\z`)
)

func normalizeField(s string) string {
	// sanitize blank characters
	s = leadingBlankRe.ReplaceAllString(s, "")
	s = trailingBlankRe.ReplaceAllString(s, "")
	return s
}

func parseAccountType(s string) AccountType {
	switch strings.ToUpper(s) {
	case "DIRECT":
		return AccountDirect
	case "RESELLER":
		return AccountReseller
	default:
		// NOTE or should be error ?
		return AccountOther
	}
}

func parseRow(row string) (Row, error) {
	if strings.Contains(row, "=") {
		// this is a variable declaration
		v, err := parseVariable(row)
		if v != nil || err != nil {
			return v, err
		}
	} else {
		// this is a record declaration
		r, err := parseRecord(row)
		if r != nil || err != nil {
			return r, err
		}
	}
	return nil, nil
}

func parseRecord(row string) (*Record, error) {
	// dropping extension field
	if idx := strings.Index(row, ";"); idx != -1 {
		row = row[0:idx]
	}

	fields := strings.Split(row, ",")

	// if the first field contains "=", then the row is for variable declaration
	if strings.Contains(fields[0], "=") {
		return nil, nil
	}

	if l := len(fields); l != 3 && l != 4 {
		return nil, fmt.Errorf("ads.txt has fields length is incorrect.: %s", row)
	}

	// otherwise the row is valid
	var r Record
	r.ExchangeDomain = strings.ToLower(normalizeField(fields[0]))
	r.PublisherAccountID = normalizeField(fields[1])
	r.AccountType = parseAccountType(normalizeField(fields[2]))
	// AuthorityID is optional
	if len(fields) >= 4 {
		r.AuthorityID = normalizeField(fields[3])
	}
	return &r, nil
}

func (r *Record) String() string {
	row := make([]string, 3, 4)
	row[0] = r.ExchangeDomain
	row[1] = r.PublisherAccountID

	switch r.AccountType {
	case AccountDirect:
		row[2] = "DIRECT"
	case AccountReseller:
		row[2] = "RESELLER"
	}

	if r.AuthorityID != "" {
		row = append(row, r.AuthorityID)
	}

	return strings.Join(row, ",")
}

// Variable is a variable row defined in iab (e.g., contact, subdomain, and such).
type Variable struct {
	Key   string
	Value string
}

func parseVariable(row string) (*Variable, error) {
	fields := strings.SplitN(row, "=", 2)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid variable row: %s", row)
	}

	return &Variable{
		Key:   normalizeField(fields[0]),
		Value: normalizeField(fields[1]),
	}, nil
}

func (v *Variable) String() string {
	return strings.Join([]string{v.Key, v.Value}, "=")
}
