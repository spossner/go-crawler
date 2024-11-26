package main

import (
	"fmt"
	"github.com/spossner/go-crawler/internal/crawler"
	"log"
	"os"
	"strconv"
)

func parseInt(src string, defaultValue int) int {
	if i, err := strconv.Atoi(src); err == nil {
		return i
	}
	return defaultValue
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("no website provided")
	}
	base := os.Args[1]
	ctx, err := crawler.NewContext(base)
	if len(os.Args) > 2 {
		ctx.SetMaxWorkers(parseInt(os.Args[2], 32))

	}

	if len(os.Args) > 3 {
		ctx.SetMaxPages(parseInt(os.Args[3], 100))

	}

	fmt.Printf("starting crawl of: %s\n", base)

	if err != nil {
		fmt.Printf("error creating crawler context: %v", err)
		os.Exit(1)
	}

	err = crawler.CrawlPage(ctx, base)
	if err != nil {
		fmt.Printf("error crawling page: %v", err)
		os.Exit(1)
	}
	ctx.WaitForWorkers()
	fmt.Println(ctx)
}
