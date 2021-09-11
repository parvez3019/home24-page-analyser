package html_parser

import (
	"golang.org/x/net/html"
	"home24-page-analyser/model"
	"io"
	"net/url"
	"strings"
)

// Parser abstraction over parser
type Parser interface {
	Parse(reader io.Reader, domain string) (model.PageAnalysisResponse, error)
}

// parser represents a struct for parser
type parser struct {
	response model.PageAnalysisResponse
	err      error
}

// NewParser creates and return a new parser object
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

// Parse takes io.reader as input and domainURL, parse the reader data and return page analysis response or error
func (p *parser) Parse(reader io.Reader, domain string) (model.PageAnalysisResponse, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return model.PageAnalysisResponse{}, err
	}

	result := p.parseHTMLVersion(node)

	var parseNode func(*html.Node)
	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			result = p.parseTitle(n).
				parseHeadings(n).
				parseLoginForm(n).
				parseLinks(n, domain)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}
	parseNode(node)

	return result.response, result.err
}

// parseHTMLVersion parse HTML version from the node
func (p *parser) parseHTMLVersion(node *html.Node) *parser {
	if len(node.FirstChild.Attr) != 0 {
		p.response.HTMLVersion = node.FirstChild.Attr[0].Val
		return p
	}
	p.response.HTMLVersion = DefaultHTMLVersion
	return p
}

// parseHeadings parse heading levels from the node
func (p *parser) parseHeadings(node *html.Node) *parser {
	if headerTagMap[node.Data] {
		p.response.HeaderCount[node.Data]++
	}
	return p
}

// parseTitle parse title from the node
func (p *parser) parseTitle(node *html.Node) *parser {
	if node.Data == TitleHTMLTagKey {
		p.response.Title = node.FirstChild.Data
	}
	return p
}

// parseLoginForm parse login form existence from the node
func (p *parser) parseLoginForm(node *html.Node) *parser {
	attrType, ok := getAttr(node.Attr, TypeHTMLTagKey)
	if ok && attrType == PasswordTypeHTMLTagKey {
		p.response.HasPasswordInputType = true
	}
	if ok && attrType == SubmitTypeHTMLTagKey {
		p.response.HasSubmitTypeInput = true
	}
	if p.response.HasSubmitTypeInput && p.response.HasPasswordInputType {
		p.response.HasLoginForm = true
	}
	return p
}

// parseLinks parse unique internal and external links from the node
func (p *parser) parseLinks(node *html.Node, domain string) *parser {
	if node.Data != AnchorHTMLTagKey {
		return p
	}
	value, found := getAttr(node.Attr, HRefKey)
	if !found || value[0] == '#' {
		// if not found or link start with # do nothing
	} else if value[0] == '/' || !startsWithHttpsScheme(value) { // if link start with / and not with any http scheme
		p.parseInternalLink(domain, value)
	} else if _, err := url.Parse(value); err == nil {
		p.response.Links.ExternalLinks.Count++
		p.response.Links.ExternalLinks.URLs = append(p.response.Links.ExternalLinks.URLs, value)
	}
	return p
}

func (p *parser) parseInternalLink(domain string, value string) {
	if err := isInternalLink(domain, value); err != nil {
		return
	}
	p.response.Links.InternalLinks.Count++
	p.response.Links.InternalLinks.URLs = append(p.response.Links.InternalLinks.URLs, value)
}

// isInternalLink verify and parse internal link
func isInternalLink(domain string, value string) error {
	urlParsed, err := url.Parse(domain)
	if err != nil {
		return err
	}
	host := urlParsed.Scheme + SchemeHostSeparator + urlParsed.Hostname()
	if _, err := url.Parse(host + value); err != nil {
		return err
	}
	return nil
}

// startsWithHttpsScheme return true if value starts with any http scheme
func startsWithHttpsScheme(value string) bool {
	return strings.Contains(value, HTTPSSchemePrefix) || strings.Contains(value, HTTPSchemePrefix)
}

// getAttr return attribute value if exists
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

const (
	TitleHTMLTagKey        = "title"
	AnchorHTMLTagKey       = "a"
	TypeHTMLTagKey         = "type"
	PasswordTypeHTMLTagKey = "password"
	SubmitTypeHTMLTagKey   = "submit"
	HRefKey                = "href"
	SchemeHostSeparator    = "://"

	HTTPSSchemePrefix = "https://"
	HTTPSchemePrefix  = "http://"
)

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
