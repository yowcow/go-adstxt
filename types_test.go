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
