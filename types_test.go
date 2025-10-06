package adstxt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRow(t *testing.T) {
	cases := []struct {
		title    string
		input    string
		expected *Row
	}{
		{
			"a direct record row",
			"hoge,fuga,DIRECT",
			&Row{
				Record: &Record{
					ExchangeDomain:     "hoge",
					PublisherAccountID: "fuga",
					AccountType:        AccountDirect,
				},
			},
		},
		{
			"a reseller record row",
			"hoge,fuga,RESELLER",
			&Row{
				Record: &Record{
					ExchangeDomain:     "hoge",
					PublisherAccountID: "fuga",
					AccountType:        AccountReseller,
				},
			},
		},
		{
			"a record row with optional AuthorityID",
			"hoge,fuga,DIRECT,123456",
			&Row{
				Record: &Record{
					ExchangeDomain:     "hoge",
					PublisherAccountID: "fuga",
					AccountType:        AccountDirect,
					AuthorityID:        "123456",
				},
			},
		},
		{
			"a record row with = in PublisherAccountID",
			"example.com,pub=123,DIRECT",
			&Row{
				Record: &Record{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "pub=123",
					AccountType:        AccountDirect,
				},
			},
		},
		{
			"a record row with = in AuthorityID",
			"example.com,12345,DIRECT,auth=xyz",
			&Row{
				Record: &Record{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "12345",
					AccountType:        AccountDirect,
					AuthorityID:        "auth=xyz",
				},
			},
		},
		{
			"a record row with extension ; should drop extension",
			"example.com,pub123,RESELLER,auth456;some extension data",
			&Row{
				Record: &Record{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "pub123",
					AccountType:        AccountReseller,
					AuthorityID:        "auth456",
				},
			},
		},
		{
			"a variable row",
			"contact=foo=bar=buz",
			&Row{
				Variable: &Variable{
					"contact",
					"foo=bar=buz",
				},
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
