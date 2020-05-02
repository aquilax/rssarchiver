package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gilliek/go-opml/opml"
	"github.com/kennygrant/sanitize"
	"github.com/mmcdole/gofeed"
)

type Archiver struct{}

func NewArchiver() *Archiver {
	return &Archiver{}
}

func (a *Archiver) Run(fileName string) error {
	doc, err := opml.NewOPMLFromFile(fileName)
	if err != nil {
		return err
	}
	links := collectLinks(doc.Body.Outlines)
	fmt.Println(links[1])

	feed, err := fetchFeed(links[1], 60)
	if err != nil {
		return err
	}
	// for _, i := range feed.Items {
	// 	println(i.Title)
	// }
	return saveFeed(links[0], feed)
}

func collectLinks(outlines []opml.Outline) []string {
	result := make([]string, 0)
	for _, o := range outlines {
		if len(o.Outlines) > 0 {
			result = append(result, collectLinks(o.Outlines)...)
		}
		if o.XMLURL != "" {
			result = append(result, o.XMLURL)
		}
	}
	return result
}

func fetchFeed(url string, timeout time.Duration) (*gofeed.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	return fp.ParseURLWithContext(url, ctx)
}

func saveFeed(url string, feed *gofeed.Feed) error {
	fileName := filepath.Join("./output", sanitize.BaseName(url))
	println(fileName)
	data, err := json.MarshalIndent(feed.Items, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}
