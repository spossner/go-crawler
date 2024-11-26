package crawler

import (
	"fmt"
	"github.com/spossner/go-crawler/internal/utils"
	"net/url"
)

func CrawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) (map[string]int, error) {
	base, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(rawCurrentURL)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != base.Host {
		// do not follow external links for now
		return pages, nil
	}
	currentUrl, err := utils.NormalizeURL(rawCurrentURL)
	if err != nil {
		return nil, err
	}

	if _, ok := pages[currentUrl]; ok {
		pages[currentUrl] += 1
		return pages, nil
	}
	pages[currentUrl] = 1

	fmt.Printf("crawling %s... ", rawCurrentURL)
	html, err := utils.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("SKIPPING")
		return pages, nil
	}
	fmt.Println("DONE")
	urls, err := utils.GetURLsFromHTML(html, rawBaseURL)
	if err != nil {
		return nil, err
	}
	for _, nextUrl := range urls {
		pages, err = CrawlPage(rawBaseURL, nextUrl, pages)
		if err != nil {
			return nil, err
		}
	}
	return pages, nil
}
