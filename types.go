package adstxt

import (
	"fmt"
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

func (v *Variable) String() string {
	return strings.Join([]string{v.Key, v.Value}, "=")
}
