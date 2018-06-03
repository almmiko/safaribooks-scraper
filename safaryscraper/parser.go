package safaryscraper

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"path/filepath"

	"golang.org/x/net/html"
)

func parseBody(content []byte) []byte {
	document, _ := html.Parse(strings.NewReader(string(content)))

	cont, err := getContent(document)
	if err != nil {
		panic(err)
	}

	c := getHtml(cont)

	return c
}

func getHtml(n *html.Node) []byte {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.Bytes()
}

func getContent(doc *html.Node) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	var navLinks []*html.Node

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {

			isLink := n.Data == "a"
			if isLink {
				for _, a := range n.Attr {
					if a.Key == "class" && (a.Val == "next nav-link" || a.Val == "prev nav-link") {
						navLinks = append(navLinks, n)
						break
					}
				}
			}

			isDiv := n.Data == "div"
			if isDiv {
				for _, a := range n.Attr {
					if a.Key == "id" && a.Val == "container" {
						b = n
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for _, link := range navLinks {
		var linkAttr []html.Attribute

		for _, a := range link.Attr {
			if a.Key == "href" {
				link := strings.TrimSuffix(a.Val, filepath.Ext(a.Val)) + ".html"
				linkAttr = append(linkAttr, html.Attribute{
					Key: a.Key,
					Val: link,
				})
			} else {
				linkAttr = append(linkAttr, html.Attribute{
					Key: a.Key,
					Val: a.Val,
				})
			}

		}

		link.Attr = linkAttr
	}

	if b != nil {
		return b, nil
	}
	return nil, errors.New("missing div#container")
}
