package util

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

//Scrapping using colly
func Scraper(url string) string {

	c := colly.NewCollector()
	title := ""

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
	})

	c.Visit(url)
	return title
}

//goquery scrapping
func RowScraper(url string) (string, error) {
	var strList []string
	var finalString string
	response, err := http.Get(url)
	if err != nil {
		return finalString, err
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return finalString, err
	}

	document.Find("a").Each(func(_ int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists {

			strList = append(strList, href)
		}
	})
	for _, v := range strList {

		finalString = fmt.Sprintf("%s \n %s", finalString, v)
	}
	return finalString, nil
}
