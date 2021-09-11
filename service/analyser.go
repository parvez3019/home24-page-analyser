package service

import (
	"home24-page-analyser/http"
	"home24-page-analyser/model"
	"home24-page-analyser/service/html_parser"
	"sync"
)

type AnalyserService interface {
	Analyse(request model.PageAnalyseRequest) (model.PageAnalysisResponse, error)
}

type analyserService struct {
	client http.Client
	parser html_parser.Parser
}

func NewAnalyzerService(httpClient http.Client, parser html_parser.Parser) AnalyserService {
	return &analyserService{
		client: httpClient,
		parser: parser,
	}
}

func (a *analyserService) Analyse(request model.PageAnalyseRequest) (model.PageAnalysisResponse, error) {
	response, err := a.client.Get(request.PageURL)
	if err != nil {
		return model.PageAnalysisResponse{}, err
	}
	pageAnalysisResponse, err := a.parser.Parse(response.Body, request.PageURL)
	if err != nil {
		return model.PageAnalysisResponse{}, err
	}

	a.verifyLinksAccessibility(&pageAnalysisResponse)
	return pageAnalysisResponse, err
}

func (a *analyserService) verifyLinksAccessibility(pageParsedResponse *model.PageAnalysisResponse) {
	inaccessibleLinksChan := make(chan string, len(pageParsedResponse.Links.ExternalLinks.URLs))
	var wg sync.WaitGroup
	for _, url := range pageParsedResponse.Links.ExternalLinks.URLs {
		wg.Add(1)
		go a.fetchURL(url, inaccessibleLinksChan, &wg)
	}
	wg.Wait()
	close(inaccessibleLinksChan)

	inaccessibleLinks := make([]string, 0)
	for inaccessibleLink := range inaccessibleLinksChan {
		inaccessibleLinks = append(inaccessibleLinks, inaccessibleLink)
	}

	pageParsedResponse.Links.InaccessibleLinks.URLs = inaccessibleLinks
	pageParsedResponse.Links.InaccessibleLinks.Count = len(inaccessibleLinks)
}

func (a *analyserService) fetchURL(url string, result chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	response, err := a.client.Get(url)
	if err != nil || response.StatusCode != 200 {
		result <- url
	}
}
