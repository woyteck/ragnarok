package indexer

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/html"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/prompter"
	"woyteck.pl/ragnarok/types"
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

func (i *Indexer) Index(document string, title string, url string) ([]uuid.UUID, error) {
	ctx := context.Background()

	doc, err := html.Parse(strings.NewReader(document))
	if err != nil {
		return nil, err
	}

	paragraphs := extractParagraphs(doc)

	memory := types.NewMemory(types.MemoryTypeWebArticle, url, document)
	i.store.Memory.InsertMemory(ctx, memory)

	insertedFragmentIds := []uuid.UUID{}
	for _, p := range paragraphs {
		content := fmt.Sprintf("%s\n%s\n\n", title, p.Title)
		for _, par := range p.Paragraphs {
			content += par + "\n"
		}

		memoryFragment := types.NewMemoryFragment(content, "", false, false, memory.ID)
		i.store.MemoryFragment.InsertMemoryFragment(ctx, memoryFragment)

		insertedFragmentIds = append(insertedFragmentIds, memoryFragment.ID)
	}

	return insertedFragmentIds, nil
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

func extractParagraphs(n *html.Node) []*Paragraph {
	results := []*Paragraph{}
	result := &Paragraph{}
	results = append(results, result)

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "p" {
				result.Paragraphs = append(result.Paragraphs, extractText(n))
			} else if n.Data == "li" {
				index := len(result.Paragraphs) - 1
				result.Paragraphs[index] = fmt.Sprintf("%s \n - %s", result.Paragraphs[index], extractText(n))
			} else if n.Data == "h1" || n.Data == "h2" || n.Data == "h3" || n.Data == "h4" || n.Data == "h5" || n.Data == "h6" {
				result = &Paragraph{
					Title: extractText(n),
				}
				results = append(results, result)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)

	paragraphs := []*Paragraph{}
	for _, p := range results {
		if len(p.Paragraphs) > 0 {
			paragraphs = append(paragraphs, p)
		}
	}

	return paragraphs
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
