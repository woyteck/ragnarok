package scraper

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

type Scrapper interface {
	ScrapPage(url string, cssSelector string) ([]string, error)
}

type CollyScraper struct {
	collector *colly.Collector
}

func NewCollyScraper(collector *colly.Collector) *CollyScraper {
	return &CollyScraper{
		collector: collector,
	}
}

func (s *CollyScraper) ScrapPage(url string, cssSelector string) ([]string, error) {
	paragraphs := []string{}

	s.collector.OnHTML(cssSelector, func(e *colly.HTMLElement) {
		text := e.Text
		text = strings.Trim(text, "\r\n")
		text = strings.Trim(text, "\t\n\v\f\r ")
		text = strings.TrimSpace(text)
		if text != "" {
			paragraphs = append(paragraphs, e.Text)
		}
	})

	err := s.collector.Visit(url)
	if err != nil {
		return nil, err
	}

	return paragraphs, nil
}

func (s *CollyScraper) GetArticle(url string, cssSelector string) (string, string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", err
	}

	title := doc.Find("title").Text()

	article := doc.Find(cssSelector).First()

	html, err := article.Html()
	if err != nil {
		return "", "", err
	}

	html = fmt.Sprintf("<div>%s</div>", html)
	html = s.cleanHtml(html, title)

	return title, html, nil
}

func (s *CollyScraper) cleanHtml(html string, title string) string {
	p := bluemonday.NewPolicy()
	p.AllowElements("p", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "li", "div", "blockquote", "q")
	p.AllowAttrs("cite").OnElements("blockquote", "q")
	html = p.Sanitize(html)

	//remove duplicated tabs
	rg := regexp.MustCompile(`[ \t]{2,}`)
	html = rg.ReplaceAllString(html, "$1")

	//strip empty tags
	html = s.removeEmptyHtmlTags(html, title)

	//remove duplicated newlines
	rg = regexp.MustCompile(`[ \n]{2,}`)
	html = rg.ReplaceAllString(html, "$1")

	//format output
	html = gohtml.Format(html)

	return html
}

func (s *CollyScraper) removeEmptyHtmlTags(h string, title string) string {
	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		panic(err)
	}

	addOrUpdateTitle(doc, title)
	removeEmptyNodes(doc)

	return renderNode(doc)
}

func isEmptyNode(n *html.Node) bool {
	if n.Type == html.ElementNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode || (c.Type == html.TextNode && strings.TrimSpace(c.Data) != "") {
				return false
			}
		}
		return true
	}
	return false
}

func removeEmptyNodes(n *html.Node) {
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		removeEmptyNodes(c)
		if isEmptyNode(c) {
			n.RemoveChild(c)
		}
		c = next
	}
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}
	return buf.String()
}

func addOrUpdateTitle(doc *html.Node, newTitle string) {
	// Find the head element
	var head *html.Node
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "html" {
			for n := c.FirstChild; n != nil; n = n.NextSibling {
				if n.Type == html.ElementNode && n.Data == "head" {
					head = n
					break
				}
			}
		}
	}

	// Create head element if it does not exist
	if head == nil {
		head = &html.Node{
			Type: html.ElementNode,
			Data: "head",
		}
		// Add the head to the html element
		doc.FirstChild.InsertBefore(head, doc.FirstChild.FirstChild)
	}

	// Find or create title element
	var title *html.Node
	for n := head.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = n
			break
		}
	}

	// Create title element if it does not exist
	if title == nil {
		title = &html.Node{
			Type: html.ElementNode,
			Data: "title",
		}
		head.AppendChild(title)
	}

	// Set the title text
	titleText := &html.Node{
		Type: html.TextNode,
		Data: newTitle,
	}

	if title.FirstChild != nil {
		title.RemoveChild(title.FirstChild)
	}
	title.AppendChild(titleText)
}
