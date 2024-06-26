package indexer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/prompter"
	"woyteck.pl/ragnarok/vectordb"
)

type Indexer struct {
	store    *db.Store
	llm      *openai.Client
	prompter *prompter.Prompter
	qdrant   *vectordb.QdrantClient
}

func NewIndexer(store *db.Store, llm *openai.Client, prompter *prompter.Prompter, qdrant *vectordb.QdrantClient) *Indexer {
	return &Indexer{
		store:    store,
		llm:      llm,
		prompter: prompter,
		qdrant:   qdrant,
	}
}

func (i *Indexer) Index(document string, title string) error {
	tempFilePathHtml := "temp/article.html"
	os.WriteFile(tempFilePathHtml, []byte(document), 0664)

	// tempFilePathTxt := "temp/article.txt"
	// text := fmt.Sprintf("%s\n\n%s", title, html2text.HTML2Text(document))
	// os.WriteFile(tempFilePathTxt, []byte(text), 0664)

	doc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		return err
	}

	paragraphs := extractTextContent(doc)
	fmt.Println(len(paragraphs))
	for _, p := range paragraphs {
		fmt.Printf("%+v\n", p)
	}

	return nil
}

func ParseHTML(r io.Reader) (*html.Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}
	return doc, nil
}

type Paragraph struct {
	Title      string
	Paragraphs []string
}

func extractTextContent(n *html.Node) []*Paragraph {
	results := []*Paragraph{}
	result := Paragraph{}

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "p" {
				result.Paragraphs = append(result.Paragraphs, extractText(n))
				if len(results) == 0 {
					results = append(results, &result)
				}
			} else if n.Data == "h1" || n.Data == "h2" || n.Data == "h3" || n.Data == "h4" || n.Data == "h5" || n.Data == "h6" {
				result = Paragraph{
					Title: extractText(n),
				}
				results = append(results, &result)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)

	return results
}

func extractText(n *html.Node) string {
	var text string
	if n.Type == html.TextNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return strings.TrimSpace(text)
}
