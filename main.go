package main

import (
	"fmt"
	"github.com/spossner/go-crawler/internal/crawler"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	base := os.Args[1]
	fmt.Printf("starting crawl of: %s\n", base)

	ctx, err := crawler.NewContext(base)
	if err != nil {
		fmt.Printf("error creating crawler context: %v", err)
		os.Exit(1)
	}

	err = crawler.CrawlPage(ctx, base)
	if err != nil {
		fmt.Printf("error crawling page: %v", err)
		os.Exit(1)
	}
	ctx.Wg.Wait()
	fmt.Println(ctx)
}
