package html_parser

import (
	"golang.org/x/net/html"
	"io"

	"home24-page-analyser/model"
)

type Parser interface {
	Parse(reader io.Reader) (model.PageAnalysisResponse, error)
}

type parser struct {
	response model.PageAnalysisResponse
	err      error
}

func NewParser() Parser {
	return &parser{}
}

func (parser *parser) Parse(reader io.Reader) (model.PageAnalysisResponse, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return model.PageAnalysisResponse{}, err
	}

	result := parser.
		parseHTMLVersion(node).
		parseHeadings(node)

	return result.response, result.err
}

func (parser *parser) parseHTMLVersion(node *html.Node) *parser {
	if len(node.FirstChild.Attr) != 0 {
		parser.response.HTMLVersion = node.FirstChild.Attr[0].Val
		return parser
	}
	parser.response.HTMLVersion = DefaultHTMLVersion
	return parser
}

func (parser *parser) parseHeadings(node *html.Node) *parser {
	headerLevelCount := make(map[string]int, 0)
	var countHeaderNode func(*html.Node)
	countHeaderNode = func(n *html.Node) {
		if n.Type == html.ElementNode && headerTagMap[n.Data] {
			headerLevelCount[n.Data]++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			countHeaderNode(c)
		}
	}
	countHeaderNode(node)
	parser.response.HeaderCount = headerLevelCount
	return parser
}

var headerTagMap = map[string]bool{
	"h1": true,
	"h2": true,
	"h3": true,
	"h4": true,
	"h5": true,
	"h6": true,
}

const DefaultHTMLVersion = "5.0"
