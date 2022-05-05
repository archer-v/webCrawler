package main

import (
	"io"

	"golang.org/x/net/html"
)

type PageData struct {
	Elements []TagItem `json:"elements"`
}

type TagItem struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}

//PageParseTagsCounter is the example of the page data parser
//it reads page data from the Reader and counts the html tags of each type
func PageParseTagsCounter(r io.Reader) (data *PageData, err error) {

	var items []TagItem
	tokens := map[string]int{}

	doc, err := html.Parse(r)
	if err != nil {
		return
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			tokens[n.Data]++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	for t, v := range tokens {
		items = append(items, TagItem{
			TagName: t,
			Count:   v,
		})
	}

	data = &PageData{
		Elements: items,
	}

	return
}
