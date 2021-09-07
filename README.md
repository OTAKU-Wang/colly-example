## 说明
使用**colly**实现对蒙古新闻网站的内容抓取，并对内容进行转换

```
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

```