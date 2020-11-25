package main

import (
	"flag"
	"fmt"
	"github.com/anaskhan96/soup"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	url string
	outputDir string
)

func main() {
	flag.StringVar(&url, "url", "", "url to crawl from")
	flag.StringVar(&outputDir, "output-dir", "", "directory to download files to")

	flag.Parse()

	if url == "" || outputDir == "" {
		panic(fmt.Errorf("you have to set --url and --output-dir"))
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	start := time.Now()

	events, err := getEvents(url)
	if err != nil {
		panic(err)
	}

	eventChan := make(chan event, len(events))
	for _, event := range events {
		eventChan <- event
	}
	close(eventChan)

	attachmentWorker := 20
	downloadWorker := 20

	attachmentChan := make(chan []attachment, attachmentWorker)

	var attachmentWg sync.WaitGroup
	attachmentWg.Add(attachmentWorker)
	for i := 0; i < attachmentWorker; i++ {
		go func() {
			defer attachmentWg.Done()
			for event := range eventChan {
				fmt.Printf("Getting attachments for event %s\n", event.name)
				eventAttachments, err := getAttachments(event)
				if err != nil {
					panic(err)
				}
				attachmentChan <- eventAttachments
			}
		}()
	}

	var downloadWg sync.WaitGroup
	downloadWg.Add(downloadWorker)
	var attachmentCounter uint64
	for i := 0; i < downloadWorker; i++ {
		go func () {
			defer downloadWg.Done()
			for attachments := range attachmentChan {
				for _, a := range attachments {
					fmt.Printf("Downloading: %s: (%s)\n", a.name, a.url)
					if err := downloadFile(a.url, path.Join(outputDir, a.name)); err != nil {
						fmt.Printf("Err: %v\n", err)
					}
					atomic.AddUint64(&attachmentCounter, 1)
				}
			}
		}()
	}

	attachmentWg.Wait()
	close(attachmentChan)

	downloadWg.Wait()

	fmt.Printf("Crawling took: %s\n", start.Sub(time.Now()))
	fmt.Printf("%d events\n", len(events))
	fmt.Printf("%d attachments\n", attachmentCounter)
}

func downloadFile(url, filepath string) error {
	// file already exists
	if _, err := os.Stat(filepath); err == nil {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

type attachment struct {
	url  string
	name string
}

func getAttachments(e event) ([]attachment, error) {
	resp, err := soup.Get(e.url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(resp)

	// get all attachment tags
	attachmentTags := doc.FindAll("a", "class", "file-uploaded")

	var ret []attachment
	for _, attachmentTag := range attachmentTags {
		if href, ok := attachmentTag.Attrs()["href"]; ok {
			ret = append(ret, attachment{
				url:  href,
				name: generateFilename(e, href),
			})
		}
	}

	return ret, nil
}

const maxFileLength = 250

var dropRegex = regexp.MustCompile(`([/?:]|%|( \(Slides Attached\)))`)

func generateFilename(e event, href string) string {
	filename := fmt.Sprintf("%s-%s", e.name, path.Base(href))
	filename = strings.ReplaceAll(filename, "%20", " ")
	filename = dropRegex.ReplaceAllString(filename, "")

	if len(filename) > maxFileLength {
		ext := filepath.Ext(filename)
		filename = filename[0:maxFileLength-len(ext)] + ext
	}
	return filename
}

type event struct {
	url  string
	name string
}

func getEvents(url string) ([]event, error) {
	resp, err := soup.Get(url)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)

	// get all event spans
	eventTags := doc.FindAll("span", "class", "event")

	var ret []event
	for _, eventTag := range eventTags {
		for _, eventChild := range eventTag.Children() {
			if href, ok := eventChild.Attrs()["href"]; ok {
				ret = append(ret, event{
					url:  fmt.Sprintf("%s/%s", url, href),
					name: eventChild.Text(),
				})
			}
		}
	}

	return ret, nil
}
