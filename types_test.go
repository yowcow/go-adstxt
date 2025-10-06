package adstxt

import "testing"

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
