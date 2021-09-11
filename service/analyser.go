package service

import (
	log "github.com/sirupsen/logrus"
	"home24-page-analyser/http"
	"home24-page-analyser/model"
	"home24-page-analyser/service/html_parser"
	httpPkg "net/http"
	"sync"
)

// AnalyserService represents an interface abstraction over analyserService
type AnalyserService interface {
	Analyse(request model.PageAnalyseRequest) (model.PageAnalysisResponse, error)
}

// analyserService represents the struct for the service
type analyserService struct {
	client http.Client
	parser html_parser.Parser
}

// NewAnalyzerService creates and returns the object for analyser service
func NewAnalyzerService(httpClient http.Client, parser html_parser.Parser) AnalyserService {
	return &analyserService{
		client: httpClient,
		parser: parser,
	}
}

// Analyse takes pageAnalyseRequest and parse the page to create page analysis response or returns error
func (a *analyserService) Analyse(request model.PageAnalyseRequest) (model.PageAnalysisResponse, error) {
	response, err := a.client.Get(request.PageURL)
	if err != nil {
		log.Errorf("Analyse failed while fetching the requested page, err : %s", err.Error())
		customErr := model.NewError(model.ErrorCodeSomethingWentWrong, err.Error(), httpPkg.StatusInternalServerError)
		return model.PageAnalysisResponse{}, customErr
	}
	pageAnalysisResponse, err := a.parser.Parse(response.Body, request.PageURL)
	if err != nil {
		log.Errorf("Analyse failed while parsing the requested page, err : %s", err.Error())
		customErr := model.NewError(model.ErrorCodeSomethingWentWrong, err.Error(), httpPkg.StatusInternalServerError)
		return model.PageAnalysisResponse{}, customErr
	}

	a.verifyLinksAccessibility(&pageAnalysisResponse)
	return pageAnalysisResponse, nil
}

// verifyLinksAccessibility makes concurrent external http call to external links to verify their accessibility
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

// fetchURL fetches the url and if it is not accessible then put it under result chan
// we assume a link is accessible it and only if is returning 200 as response
// Improvement Scope - we can also consider 3xx as valid response statues but not implemented yet
func (a *analyserService) fetchURL(url string, inaccessibleLinksChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	response, err := a.client.Get(url)
	if err != nil || response.StatusCode != 200 {
		inaccessibleLinksChan <- url
	}
}
