package utils

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/url"
	"strings"
)

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	urls := make([]string, 0)

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					u, err := url.Parse(a.Val)
					if err != nil {
						return nil, err
					}
					foundUrl := base.ResolveReference(u)
					urls = append(urls, foundUrl.String())
					break
				}
			}
		}
	}
	return urls, nil
}
