package main

import (
	"fmt"
	"github.com/spossner/go-crawler/internal/utils"
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

	html, err := utils.GetHTML(base)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(html)
}
