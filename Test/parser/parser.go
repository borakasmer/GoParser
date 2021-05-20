package parser

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

var articleList map[string]string

func ParseWeb() map[string]string {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Get("https://www.borakasmer.com")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			data := doc.Find(".entry-title")
			articleList = make(map[string]string, data.Length())
			data.Each(func(i int, s *goquery.Selection) {
				title := s.Find("a").Text()
				url, _ := s.Find("a").Attr("href")
				articleList[title] = url
			})
			return articleList
		}
	} else {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return articleList
}
