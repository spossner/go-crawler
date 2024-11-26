package crawler

import (
	"fmt"
	"github.com/spossner/go-crawler/internal/utils"
	"net/url"
	"sync"
)

type Context struct {
	pages      map[string]int
	baseURL    *url.URL
	mu         *sync.RWMutex
	Wg         *sync.WaitGroup
	maxWorkers chan struct{}
}

func NewContext(rawBaseURL string) (*Context, error) {
	base, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}
	return &Context{
		pages:      make(map[string]int),
		baseURL:    base,
		mu:         &sync.RWMutex{},
		Wg:         &sync.WaitGroup{},
		maxWorkers: make(chan struct{}, 10),
	}, nil
}

func (c *Context) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.pages[key]; ok {
		return true
	}
	return false
}

func (c *Context) Set(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pages[key] = value
}

func (c *Context) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pages[key] += 1
}

func (c *Context) String() string {
	return fmt.Sprintf("%v", c.pages)
}

func CrawlPage(ctx *Context, rawCurrentURL string) error {
	u, err := url.Parse(rawCurrentURL)
	if err != nil {
		return err
	}
	if u.Host != "" && u.Host != ctx.baseURL.Host {
		// do not follow external links for now
		return nil
	}
	currentUrl, err := utils.NormalizeURL(rawCurrentURL)
	if err != nil {
		return err
	}

	if ctx.Has(currentUrl) {
		ctx.Inc(currentUrl)
		return nil
	}
	ctx.Set(currentUrl, 1)

	html, err := utils.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("skipping %s\n", rawCurrentURL)
		return nil
	}
	fmt.Printf("crawling %s\n", rawCurrentURL)
	urls, err := utils.GetURLsFromHTML(html, ctx.baseURL)
	if err != nil {
		return err
	}
	for _, nextUrl := range urls {
		ctx.Wg.Add(1)
		go func() {
			defer ctx.Wg.Done()
			ctx.maxWorkers <- struct{}{}
			CrawlPage(ctx, nextUrl)
			<-ctx.maxWorkers
		}()
	}
	return nil
}
