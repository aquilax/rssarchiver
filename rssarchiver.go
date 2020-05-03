package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/gilliek/go-opml/opml"
	"github.com/kennygrant/sanitize"
	"github.com/mmcdole/gofeed"
)

type Archiver struct{}

func NewArchiver() *Archiver {
	return &Archiver{}
}

func (a *Archiver) UpdateFromOPML(fileName string) error {
	doc, err := opml.NewOPMLFromFile(fileName)
	if err != nil {
		return err
	}
	links := collectLinks(doc.Body.Outlines)

	var workers = runtime.NumCPU()
	if workers > 1 {
		workers--
	}
	wp := workerpool.New(workers)

	total := len(links)

	for i, l := range links {
		l := l
		i := i + 1
		wp.Submit(func() {
			log.Printf("[%d/%d] Procesing %v", i, total, l)
			feed, err := fetchFeed(l, 10)
			if err != nil {
				log.Printf("[%d/%d] Error downloading feed %s : %v", i, total, l, err)
				return
			}
			err = saveFeed(l, feed)
			if err != nil {
				log.Printf("[%d/%d] Error saving feed %s : %v", i, total, l, err)
			}
			log.Printf("[%d/%d] Done with %s", i, total, l)
		})
	}
	wp.StopWait()
	return nil
}

// GenerateSummary generates feeds summary for a specified date
func (a *Archiver) GenerateSummary(date string) error {
	return fmt.Errorf("Not implemented")
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
	var err error
	var data []byte
	fileName := filepath.Join("./output", sanitize.BaseName(url)) + ".json"

	var downloadedFeed gofeed.Feed
	_, err = os.Stat(fileName)
	if !os.IsNotExist(err) {
		data, err = ioutil.ReadFile(fileName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &downloadedFeed)
		if err != nil {
			return err
		}
		appendFeed(feed, downloadedFeed)
	}
	data, err = json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

func appendFeed(feed *gofeed.Feed, downloadedFeed gofeed.Feed) {
	newIds := make(map[string]interface{})
	for _, i := range feed.Items {
		newIds[i.GUID] = nil
	}
	for _, i := range downloadedFeed.Items {
		if _, ok := newIds[i.GUID]; !ok {
			feed.Items = append(feed.Items, i)
		}
	}
}
