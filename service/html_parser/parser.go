package html_parser

import (
	"golang.org/x/net/html"
	"home24-page-analyser/model"
	"io"
	"net/url"
	"strings"
)

type Parser interface {
	Parse(reader io.Reader, domain string) (model.PageAnalysisResponse, error)
}

type parser struct {
	response model.PageAnalysisResponse
	err      error
}

func NewParser() Parser {
	return &parser{
		response: model.PageAnalysisResponse{
			HeaderCount: make(map[string]int, 0),
			Links: model.LinksResponse{
				InternalLinks:     model.LinkCountResponse{URLs: make([]string, 0)},
				ExternalLinks:     model.LinkCountResponse{URLs: make([]string, 0)},
				InaccessibleLinks: model.LinkCountResponse{URLs: make([]string, 0)},
			},
		},
	}
}

func (parser *parser) Parse(reader io.Reader, domain string) (model.PageAnalysisResponse, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return model.PageAnalysisResponse{}, err
	}

	result := parser.
		parseTitle(node).
		parseHTMLVersion(node).
		parseHeadings(node).
		parseLoginForm(node).
		parseLinks(node, domain)

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

func (parser *parser) parseTitle(node *html.Node) *parser {
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			parser.response.Title = n.FirstChild.Data
			return
		}
		for childNode := n.FirstChild; childNode != nil; childNode = childNode.NextSibling {
			findTitle(childNode)
		}
	}
	findTitle(node)
	return parser
}

func (parser *parser) parseLoginForm(doc *html.Node) *parser {
	var hasPasswordInputType, hasSubmitTypeInput bool
	var findLoginForm func(*html.Node)
	findLoginForm = func(n *html.Node) {
		if n.Type == html.ElementNode {
			attrType, ok := getAttr(n.Attr, "type")
			if ok && attrType == "password" {
				hasPasswordInputType = true
			}
			if ok && attrType == "submit" {
				hasSubmitTypeInput = true
			}
			if hasSubmitTypeInput && hasPasswordInputType {
				parser.response.HasLoginForm = true
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLoginForm(c)
		}
	}
	findLoginForm(doc)
	return parser
}

func (parser *parser) parseLinks(doc *html.Node, domain string) *parser {
	tags := make([]AnchorTag, 0)
	var parseNodeLink func(*html.Node)
	parseNodeLink = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			value, validLink := getAttr(n.Attr, "href")
			if !validLink {
				return
			}
			if value[0] == '/' || !startsWithHttpsScheme(value) {
				urlParsed, err := url.Parse(domain)
				if err != nil {
					return
				}

				host := urlParsed.Scheme + "://" + urlParsed.Hostname()
				if _, err = url.Parse(host + value); err == nil {
					tags = append(tags, AnchorTag{
						Url:        value,
						IsExternal: false,
					})
				}
			} else if _, err := url.Parse(value); err == nil {
				tags = append(tags, AnchorTag{
					Url:        value,
					IsExternal: true,
				})
				return
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNodeLink(c)
		}
	}
	parseNodeLink(doc)

	for _, tag := range tags {
		if tag.IsExternal {
			parser.response.Links.ExternalLinks.Count++
			parser.response.Links.ExternalLinks.URLs = append(parser.response.Links.ExternalLinks.URLs, tag.Url)
		} else {
			parser.response.Links.InternalLinks.Count++
			parser.response.Links.InternalLinks.URLs = append(parser.response.Links.InternalLinks.URLs, tag.Url)
		}
	}
	return parser
}

func startsWithHttpsScheme(value string) bool {
	return strings.Contains(value, "https://") || strings.Contains(value, "http://")
}

func getAttr(attrs []html.Attribute, attrName string) (string, bool) {
	var attrVal string
	var found bool
	for _, attr := range attrs {
		if attr.Key == attrName {
			attrVal = attr.Val
			found = true
		}
	}
	return attrVal, found
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

type AnchorTag struct {
	Url        string
	IsExternal bool
}
