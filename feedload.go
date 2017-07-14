package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/asaskevich/govalidator"
)

const WORKERS = 8

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("No target rss url specified")
	}

	rssUrl := os.Args[1]
	if !govalidator.IsURL(rssUrl) {
		log.Fatal("Your argument is no valid url")
	}

	download(rssUrl)
	fmt.Println("Finished")
}

func download(url string) {
	progress := mpb.New()

	//start workers
	jobs := make(chan gofeed.Item, 2)
	done := make(chan struct{}, 2)
	for w := 0; w < WORKERS; w++ {
		go worker(jobs, done, progress)
	}

	//read rss
	feed, err := gofeed.NewParser().ParseURL(url)
	if err != nil {
		log.Fatal("Cannot parse feed source", err)
	}

	//count all episodes to make an progress bar
	totalEpisodes := len(feed.Items)
	go progressCounter(progress, totalEpisodes, done)

	//submit all rss entries to the workers
	fmt.Println("Downloading all episodes from", feed.Title)
	for _, item := range feed.Items {
		jobs <- *item
	}

	//Wait and clean up
	close(jobs)
	close(done)
	progress.Stop()
}

func progressCounter(progress *mpb.Progress, totalEpisodes int, done <-chan struct{}) {
	bar := createBar(progress, "All", 0, int64(totalEpisodes))

	for range done {
		bar.Incr(1)
	}
}

func worker(jobs <-chan gofeed.Item, done chan<- struct{}, progress *mpb.Progress) {
	for item := range jobs {
		enclosures := item.Enclosures
		downloadSource := enclosures[0].URL

		fileExt := extractFileExt(downloadSource)

		downloadFile(item.Title, fileExt, downloadSource, progress)
		done <- struct{}{}
	}
}

func extractFileExt(downloadSource string) string {
	//extract the file extension from the url
	split := strings.Split(downloadSource, ".")
	return split[len(split) - 1]
}

func downloadFile(title string, ext string, url string, progress *mpb.Progress) {
	//create output file
	out, err := os.Create(title + "." + ext)
	defer out.Close()
	if err != nil {
		//cancel processing of this file
		log.Println("Failed to create file", err)
		return
	}

	//connect with http
	resp, err := http.Get(url)
	defer resp.Body.Close()

	//add progress with the maximum file size
	bar := createBar(progress, title, decor.Unit_KiB, resp.ContentLength)
	defer progress.RemoveBar(bar)

	//write contents to disk
	_, err = io.Copy(out, bar.ProxyReader(resp.Body))
	if err != nil {
		log.Println("Failed to download file", err)
	}
}

func createBar(progress *mpb.Progress, name string, counterUnit decor.Units, total int64) *mpb.Bar {
	return progress.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(name, 0, decor.DwidthSync),
			decor.Counters("%3s/%3s", counterUnit, 18, 0),
		),
		mpb.AppendDecorators(
			decor.ETA(3, 0),
			decor.Percentage(5, 0),
		))
}
