package adstxt

import (
	"fmt"
	"regexp"
	"strings"
)

// Row represents a parsed ads.txt row.
// Exactly one of Record or Variable will be non-nil when parsing succeeds.
// Row is nil when an ads.txt row is not a valid row.
type Row struct {
	Record   *Record
	Variable *Variable
}

type AccountType int

const (
	AccountDirect AccountType = iota
	AccountReseller
	AccountOther
)
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

func parseRow(row string) (*Row, error) {
	// Check if "=" exists before the first ","
	commaIndex := strings.Index(row, ",")
	equalsIndex := strings.Index(row, "=")

	if equalsIndex != -1 && (commaIndex == -1 || equalsIndex < commaIndex) {
		// this is a variable declaration
		v, err := parseVariable(row)
		if err != nil {
			return nil, err
		}
		if v != nil {
			return &Row{Variable: v}, nil
		}
	} else {
		// this is a record declaration
		r, err := parseRecord(row)
		if err != nil {
			return nil, err
		}
		if r != nil {
			return &Row{Record: r}, nil
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

	if l := len(fields); l != 3 && l != 4 {
		return nil, fmt.Errorf("ads.txt field length is incorrect: %s", row)
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

// Variable represents a variable declaration in the ads.txt specification.
// Variables are used to provide additional metadata or configuration, such as contact information or subdomain delegation.
// For example, a variable row in ads.txt might look like:
//
//	contact=ads@example.com
//	subdomain=ads.subdomain.com
//
// Key is the variable name (e.g., "contact", "subdomain").
// Value is the value assigned to the variable (e.g., "ads@example.com", "ads.subdomain.com").
type Variable struct {
	// Key is the variable name in the ads.txt variable declaration.
	// Example: "contact", "subdomain"
	Key string
	// Value is the value assigned to the variable in the ads.txt variable declaration.
	// Example: "ads@example.com", "ads.subdomain.com"
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
