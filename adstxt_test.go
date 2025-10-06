package adstxt_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/suzuken/go-adstxt"
)

func TestParse(t *testing.T) {
	cases := []struct {
		txt      string
		expected []adstxt.Record
	}{
		{
			txt: `example.com,1,DIRECT`,
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
			},
		},
		{
			txt: "example.com,1,DIRECT\nexample.org,2,RESELLER",
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
				{
					ExchangeDomain:     "example.org",
					PublisherAccountID: "2",
					AccountType:        adstxt.AccountReseller,
				},
			},
		},
		{
			txt: "\n\nEXAMPLE.COM, 1, direct, TAG ID1; COMMENT1\n\nEXAMPLE.ORG , \t2 , \treseller , \tTAG ID2 ; \tCOMMENT2\n\nfoo=bar",
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
					AuthorityID:        "TAG ID1",
				},
				{
					ExchangeDomain:     "example.org",
					PublisherAccountID: "2",
					AccountType:        adstxt.AccountReseller,
					AuthorityID:        "TAG ID2",
				},
			},
		},
		{
			txt: "# comment.out,comment-publisher,DIRECT\nexample.com,1,DIRECT",
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
			},
		},
		{
			txt: "example.com,1,DIRECT# trailing comment\nexample.com,2,RESELLER###trailing comment",
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "2",
					AccountType:        adstxt.AccountReseller,
				},
			},
		},
		{
			txt: "example.com,1,DIRECT\n      \nexample.com,2,RESELLER",
			expected: []adstxt.Record{
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "1",
					AccountType:        adstxt.AccountDirect,
				},
				{
					ExchangeDomain:     "example.com",
					PublisherAccountID: "2",
					AccountType:        adstxt.AccountReseller,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.txt, func(t *testing.T) {
			record, err := adstxt.Parse(strings.NewReader(c.txt))
			if err != nil {
				t.Errorf("parse ads.txt failed: %s", err)
			}
			if d := cmp.Diff(c.expected, record); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestParseRows(t *testing.T) {
	cases := []struct {
		txt      string
		expected []adstxt.Row
	}{
		{
			txt: "example.com,1,DIRECT",
			expected: []adstxt.Row{
				{
					Record: &adstxt.Record{
						ExchangeDomain:     "example.com",
						PublisherAccountID: "1",
						AccountType:        adstxt.AccountDirect,
					},
				},
			},
		},
		{
			txt: "example.com,1,DIRECT\nfoo=bar",
			expected: []adstxt.Row{
				{
					Record: &adstxt.Record{
						ExchangeDomain:     "example.com",
						PublisherAccountID: "1",
						AccountType:        adstxt.AccountDirect,
					},
				},
				{
					Record:   nil,
					Variable: &adstxt.Variable{Key: "foo", Value: "bar"},
				},
			},
		},
		{
			txt: "# comment\nexample.org,2,RESELLER\nbaz=qux",
			expected: []adstxt.Row{
				{
					Record: &adstxt.Record{
						ExchangeDomain:     "example.org",
						PublisherAccountID: "2",
						AccountType:        adstxt.AccountReseller,
					},
					Variable: nil,
				},
				{
					Record:   nil,
					Variable: &adstxt.Variable{Key: "baz", Value: "qux"},
				},
			},
		},
		{
			txt: "   \nexample.com,1,DIRECT\n# a comment line\nfoo=bar\nexample.org,2,RESELLER",
			expected: []adstxt.Row{
				{
					Record: &adstxt.Record{
						ExchangeDomain:     "example.com",
						PublisherAccountID: "1",
						AccountType:        adstxt.AccountDirect,
					},
					Variable: nil,
				},
				{
					Record:   nil,
					Variable: &adstxt.Variable{Key: "foo", Value: "bar"},
				},
				{
					Record: &adstxt.Record{
						ExchangeDomain:     "example.org",
						PublisherAccountID: "2",
						AccountType:        adstxt.AccountReseller,
					},
					Variable: nil,
				},
			},
		},
		{
			txt: "key=value\nexample.com,1,DIRECT\nkey2=value2",
			expected: []adstxt.Row{
				{
					Record:   nil,
					Variable: &adstxt.Variable{Key: "key", Value: "value"},
				},
				{
					Record: &adstxt.Record{
						ExchangeDomain:     "example.com",
						PublisherAccountID: "1",
						AccountType:        adstxt.AccountDirect,
					},
					Variable: nil,
				},
				{
					Record:   nil,
					Variable: &adstxt.Variable{Key: "key2", Value: "value2"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.txt, func(t *testing.T) {
			rows, err := adstxt.ParseRows(strings.NewReader(c.txt))
			if err != nil {
				t.Errorf("parse rows failed: %s", err)
			}
			if d := cmp.Diff(c.expected, rows); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}
