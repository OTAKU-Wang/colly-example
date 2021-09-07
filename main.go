package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	cuitil "example.com/crawl/pkg"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(colly.AllowedDomains("mgyxw.cn", "www.mgyxw.cn"))
	q, _ := queue.New(
		1, &queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	c.OnHTML("div[id=mk_nr]", func(e *colly.HTMLElement) {
		log.Println("content found", e.Request.URL)
		fName := fmt.Sprintf("%d.txt", e.Request.ID)
		file, err := os.Create(fName)
		if err != nil {
			log.Fatalf("Cannot create file %q: %s\n", fName, err)
			return
		}
		defer file.Close()
		content := e.ChildText("span.mkh_ctt")
		file.WriteString(content)
		log.Println(content)
	})

	for i := 1; i <= 1; i++ {
		// Add URLs to the queue
		timeUnix := time.Now().Unix()
		urlTemp := "http://api.mgyxw.cn/mdls/am/amList.ashx?call=ArticleMore&st=_H&mid=9376&ct=mk&mh=390&cnt=39&pn=%d&_=%d"
		s := fmt.Sprintf(urlTemp, i, timeUnix)
		if s, err := cuitil.GetRequestUrl(s, handleBody); err == nil {
			for _, v := range s {
				q.AddURL(v)
			}
		} else {
			fmt.Println(err)
		}
	}
	// Consume URLs
	q.Run(c)

}

type ArticleMore struct {
	Mk_Contents []Content `json:"mk_Contents"`
}
type Content struct {
	Mk_STT string `json:"mk_STT"`
	Mk_URL string `json:"mk_URL"`
}

func handleBody(body string) ([]string, error) {
	articleMore := ArticleMore{}
	urlPreFix := "http://www.mgyxw.cn"
	var urls []string
	jsonString := body[12 : len(body)-1]
	if err := json.Unmarshal([]byte(jsonString), &articleMore); err == nil {
		c := articleMore.Mk_Contents
		for _, v := range c {
			fmt.Println(v.Mk_STT)
			urls = append(urls, urlPreFix+v.Mk_URL)
		}
	} else {
		return nil, err
	}
	return urls, nil
}
