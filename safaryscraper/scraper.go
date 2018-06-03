package safaryscraper

import (
	"io"
	"log"
	"net/http"
	"os"

	"path/filepath"

	"strings"

	"fmt"

	"strconv"

	"github.com/gocolly/colly"
)

type bookScrapper struct {
	Config *Config
}

func NewBookScrapper(config *Config) *bookScrapper {
	return &bookScrapper{
		Config: config,
	}
}

func (bs *bookScrapper) GetHtmlPages() {

	var pageCount int

	c := colly.NewCollector()

	cookies := newCookiesList(bs.Config.RawCookies)

	c.SetCookies(bs.Config.Url, cookies)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		writeHtml(e.Request.URL.Path, e.Response.Body)
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		saveImage(e.Request.AbsoluteURL(src), src)
	})

	c.OnHTML(".t-sbo-next.sbo-next.sbo-nav-top .next.nav-link", func(e *colly.HTMLElement) {
		pageCount++
		fmt.Println("Processing page " + strconv.Itoa(pageCount))
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit(bs.Config.Url)
}

func createDir(filePath string) {
	dirPath, _ := filepath.Split(filePath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		panic(err)
	}
}

func saveImage(url string, path string) {

	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	filePath, err := filepath.Abs("../BooksScrapper/html/" + path)
	if err != nil {
		panic(err)
	}

	createDir(filePath)

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
}

func writeHtml(path string, content []byte) {

	fileName := strings.TrimSuffix(path, filepath.Ext(path)) + ".html"

	filePath, err := filepath.Abs("../BooksScrapper/html/" + fileName)
	if err != nil {
		panic(err)
	}

	createDir(filePath)

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	c := parseBody(content)
	styles := getStyles()
	withStyles := append(c, styles...)

	if _, err := file.Write(withStyles); err != nil {
		panic(err)
	}
}
