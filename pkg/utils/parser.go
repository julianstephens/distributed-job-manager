package utils

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/html"
)

type Parser struct {
	SanitizedInput string
	Result         []CodeBlock
}

type CodeBlock struct {
	Language string
	Content  string
}

type DocumentData struct {
	CodeBlocks []CodeBlock
}

func (p *Parser) Parse(unsafeInput string) error {
	unsafe := blackfriday.Run([]byte(unsafeInput))
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	res := policy.SanitizeBytes(unsafe)

	doc, err := html.Parse(bytes.NewReader(res))
	if err != nil {
		return err
	}

	data := DocumentData{}
	p.ExtractCodeBlocks(doc, &data)

	p.SanitizedInput = string(res)
	p.Result = data.CodeBlocks

	return nil
}

func (p *Parser) ExtractCodeBlocks(n *html.Node, data *DocumentData) {
	if n.Type == html.ElementNode && n.Data == "pre" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "code" {
				language := ""
				for _, attr := range c.Attr {
					if attr.Key == "class" && strings.HasPrefix(attr.Val, "language-") {
						language = strings.TrimPrefix(attr.Val, "language-")
						break
					}
				}
				content := c.FirstChild.Data
				data.CodeBlocks = append(data.CodeBlocks, CodeBlock{
					Language: language,
					Content:  strings.ReplaceAll(content, "'", "\""),
				})
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.ExtractCodeBlocks(c, data)
	}
}
