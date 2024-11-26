package crawler

import (
	"fmt"
	"github.com/spossner/go-crawler/internal/utils"
	"golang.org/x/exp/constraints"
	"net/url"
	"sort"
	"strings"
	"sync"
)

type Context struct {
	pages      map[string]int
	baseURL    *url.URL
	mu         *sync.RWMutex
	wg         *sync.WaitGroup
	maxWorkers chan struct{}
	maxPages   int
}

type ValueKeyPair[K comparable, T constraints.Ordered] struct {
	Key   K
	Value T
}
type SortedByValue[K comparable, T constraints.Ordered] []ValueKeyPair[K, T]

func (s SortedByValue[K, T]) Len() int           { return len(s) }
func (s SortedByValue[K, T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortedByValue[K, T]) Less(i, j int) bool { return s[i].Value < s[j].Value }

func (s SortedByValue[K, T]) Sort(asc bool) {
	if asc {
		sort.Sort(s)
	} else {
		sort.Sort(sort.Reverse(s))
	}
}

func AsSortedByValue[K comparable, T constraints.Ordered](m map[K]T, asc bool) (s SortedByValue[K, T]) {
	for k, v := range m {
		s = append(s, ValueKeyPair[K, T]{k, v})
	}

	s.Sort(asc)
	return
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
		wg:         &sync.WaitGroup{},
		maxWorkers: make(chan struct{}, 32),
		maxPages:   100,
	}, nil
}

func (c *Context) SetMaxWorkers(maxWorkers int) {
	c.maxWorkers = make(chan struct{}, maxWorkers)
}

func (c *Context) SetMaxPages(maxPages int) {
	c.maxPages = maxPages
}

func (c *Context) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.pages[key]; ok {
		return true
	}
	return false
}

func (c *Context) MaxPagesReached() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.pages) >= c.maxPages
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

func (c *Context) WaitForWorkers() {
	c.wg.Wait()
}

func (c *Context) String() string {
	buf := []string{
		"=============================",
		fmt.Sprintf("REPORT for %s", c.baseURL.String()),
		"=============================",
	}
	sorted := AsSortedByValue(c.pages, false)
	for _, k := range sorted {
		buf = append(buf, fmt.Sprintf("Found %d internal links to %v", k.Value, k.Key))
	}
	return strings.Join(buf, "\n")
}

func CrawlPage(ctx *Context, rawCurrentURL string) error {
	if ctx.MaxPagesReached() {
		return nil
	}
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
		ctx.wg.Add(1)
		go func() {
			defer ctx.wg.Done()
			ctx.maxWorkers <- struct{}{}
			CrawlPage(ctx, nextUrl)
			<-ctx.maxWorkers
		}()
	}
	return nil
}
