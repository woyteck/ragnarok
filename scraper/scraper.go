package scraper

import (
	"strings"

	"github.com/gocolly/colly"
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
