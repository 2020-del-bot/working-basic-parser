package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	url := "http://quotes.toscrape.com"

	fmt.Print("\n/==========================================================================================/\n\n")

	c := colly.NewCollector(
		colly.AllowedDomains("quotes.toscrape.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("!!!Visiting %s\n", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("The site coulddnt get reached", err, "site", r.Request.URL)
	})

	c.OnHTML(".next a", func(h *colly.HTMLElement) {
		nextpage := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(nextpage)
	})

	c.OnHTML(".quote", func(h *colly.HTMLElement) {

		text := h.ChildText(".text")

		author := h.ChildText(".author")

		var tags []string

		h.ForEach(".tag", func(_ int, tag *colly.HTMLElement) {
			tags = append(tags, tag.Text)
		})

		fmt.Printf("!!!TEXT %s\n", text)
		fmt.Printf("!!!AUTHOR %s\n", author)
		fmt.Printf("!!!TAG %s\n\n", strings.Join(tags, ", "))

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Print("/==========================================================================================/\n\n")

	})

	err := c.Visit(url)

	if err != nil {
		fmt.Println("Cant connect")
	}
}
