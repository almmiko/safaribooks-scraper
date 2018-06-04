package safaryscraper

import (
	"io"
	"log"
	"net/http"
	"os"

	"path/filepath"

	"strings"

	"sync"

	"fmt"

	"strconv"

	"github.com/gocolly/colly"
)

type bookScrapper struct {
	Config     *Config
	BookStyles []byte
	once       sync.Once
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
		bs.once.Do(func() {
			styles := getStyles(e)
			bs.BookStyles = append(bs.BookStyles, styles...)
		})

		bs.writeHtml(e.Request.URL.Path, e.Response.Body, e)
	})

	c.OnHTML("link", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fetchStyles(e.Request.AbsoluteURL(link), link)
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

func (bs *bookScrapper) writeHtml(path string, content []byte, e *colly.HTMLElement) {

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
	withStyles := append(c, bs.BookStyles...)

	if _, err := file.Write(withStyles); err != nil {
		panic(err)
	}
}

func getStyles(html *colly.HTMLElement) []byte {
	var styles []byte

	s := html.DOM.Find("style")

	for _, n := range s.Nodes {
		styleHtml := getHtml(n)
		styles = append(styles, styleHtml...)
	}

	return styles
}

func fetchStyles(url string, path string) {

	filePath, err := filepath.Abs("../BooksScrapper/html/" + path)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(filePath); err == nil {
		return
	}

	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

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
