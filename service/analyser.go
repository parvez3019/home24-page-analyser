package service

import (
	"home24-page-analyser/http"
	"home24-page-analyser/model"
	"home24-page-analyser/service/html_parser"
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
	return a.parser.Parse(response.Body)
}
