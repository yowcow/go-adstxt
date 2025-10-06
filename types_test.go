package adstxt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRow(t *testing.T) {
	cases := []struct {
		title    string
		input    string
		expected Row
	}{
		{
			"a direct record row",
			"hoge,fuga,DIRECT",
			&Record{
				ExchangeDomain:     "hoge",
				PublisherAccountID: "fuga",
				AccountType:        AccountDirect,
			},
		},
		{
			"a reseller record row",
			"hoge,fuga,RESELLER",
			&Record{
				ExchangeDomain:     "hoge",
				PublisherAccountID: "fuga",
				AccountType:        AccountReseller,
			},
		},
		{
			"a record row with optional AuthorityID",
			"hoge,fuga,DIRECT,123456",
			&Record{
				ExchangeDomain:     "hoge",
				PublisherAccountID: "fuga",
				AccountType:        AccountDirect,
				AuthorityID:        "123456",
			},
		},
		{
			"a variable row",
			"contact=foo=bar=buz",
			&Variable{
				"contact",
				"foo=bar=buz",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			r, err := parseRow(c.input)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			if d := cmp.Diff(c.expected, r); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestRow(t *testing.T) {
	cases := []struct {
		title     string
		item      Row
		expectedd string
	}{
		{
			"DIRECT record",
			&Record{
				ExchangeDomain:     "example.com",
				PublisherAccountID: "123",
				AccountType:        AccountDirect,
			},
			"example.com,123,DIRECT",
		},
		{
			"RESELLER record",
			&Record{
				ExchangeDomain:     "example.com",
				PublisherAccountID: "123",
				AccountType:        AccountReseller,
			},
			"example.com,123,RESELLER",
		},
		{
			"a record with optional AuthorityID",
			&Record{
				ExchangeDomain:     "example.com",
				PublisherAccountID: "123",
				AccountType:        AccountDirect,
				AuthorityID:        "TAGID123",
			},
			"example.com,123,DIRECT,TAGID123",
		},
		{
			"a variable",
			&Variable{
				"contact",
				"foo@example.com",
			},
			"contact=foo@example.com",
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			actual := c.item.String()
			if c.expectedd != actual {
				t.Errorf("expected %q, but got %q", c.expectedd, actual)
			}
		})
	}
}
