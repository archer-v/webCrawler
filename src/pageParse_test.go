package main

import (
	"strconv"
	"strings"
	"testing"
)

type pageParseTestCase struct {
	html 		string
	items       []TagItem
}

func TestPageParse(t *testing.T) {

	cases := []pageParseTestCase{
		{`<html><head></head><body><p>Test:</p><ul><li><a href="one">One</a></li></ul><p><a href="two">Two</a></p><div><p>The end</p></div></body>`,
			[]TagItem{
			{"p", 3},
			{"ul", 1},
			{"li", 1},
			{"a", 2},
			{"div", 1},
			{"html", 1},
			{"head", 1},
			{"body", 1},
			},
		},
		{``,
			[]TagItem{
				{"html", 1},
				{"head", 1},
				{"body", 1},
			},
		},
	}

	for i, c := range cases {
		t.Run("Case " + strconv.Itoa(i), func(t *testing.T) {
			data, err := PageParseTagsCounter(strings.NewReader(c.html))
			if err != nil {
				t.Fatalf("Expected error %v", err)
			}
			if len(data.Elements) != len(c.items) {
				t.Fatalf("Expected %v tags, got %v", len(c.items), len(data.Elements))
			}
			for _, e := range data.Elements {
				found := false
				for _, ce := range c.items {
					if ce.TagName == e.TagName && ce.Count == e.Count {
						found = true
						break
					}
				}
				if ! found {
					t.Fatalf("Count of tag '%v' is wrong", e.TagName)
				}
			}
		})
	}
}