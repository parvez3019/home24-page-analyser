package service

import "home24-page-analyser/model"

type AnalyserService interface {
	Analyse(request interface{}) (model.PageAnalysisResponse, error)
}

type analyserService struct {
}

func NewAnalyzerService() AnalyserService {
	return &analyserService{}
}

func (a *analyserService) Analyse(request interface{}) (model.PageAnalysisResponse, error) {
	return model.PageAnalysisResponse{}, nil
}
