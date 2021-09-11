package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"home24-page-analyser/model"
	"home24-page-analyser/service"
	"home24-page-analyser/service/mocks"
	"home24-page-analyser/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_pageAnalyserHandler_Analyse(t *testing.T) {
	tests := []struct {
		name               string
		context            *gin.Context
		analyserService    service.AnalyserService
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "should return 200 with no error for a valid page request",
			analyserService: func() service.AnalyserService {
				mockedAnalyserService := &mocks.AnalyserService{}
				mockedAnalyserService.On("Analyse", model.PageAnalyseRequest{PageURL: "test.com"}).
					Return(model.PageAnalysisResponse{HTMLVersion: "Version 5", Title: "Title"}, nil)
				return mockedAnalyserService
			}(),
			context:            getMockContextWithValidRequestBody(),
			expectedStatusCode: 200,
			expectedResponse:   utils.LoadJSONFromPath("/home24-page-analyser/handler/test_data/valid_page_analysis_response.json"),
		},
		{
			name:               "should return 400 with error for an invalid page request",
			analyserService:    nil,
			context:            getMockContext(),
			expectedStatusCode: 400,
			expectedResponse:   `{"error": "INVALID_REQUEST","message": "invalid request","statusCode": 400}`,
		},
		{
			name: "should return error when analyser service return custom error",
			analyserService: func() service.AnalyserService {
				mockedAnalyserService := &mocks.AnalyserService{}
				mockedAnalyserService.On("Analyse", model.PageAnalyseRequest{PageURL: "test.com"}).
					Return(model.PageAnalysisResponse{}, model.NewError(model.ErrorCodeSomethingWentWrong, "Error", 500))
				return mockedAnalyserService
			}(),
			context:            getMockContextWithValidRequestBody(),
			expectedStatusCode: 500,
			expectedResponse:   `{"error":"SOMETHING_WENT_WRONG","message":"Error","statusCode":500}`,
		},
		{
			name: "should return error when analyser service return go basic error",
			analyserService: func() service.AnalyserService {
				mockedAnalyserService := &mocks.AnalyserService{}
				mockedAnalyserService.On("Analyse", model.PageAnalyseRequest{PageURL: "test.com"}).
					Return(model.PageAnalysisResponse{}, errors.New("someErr"))
				return mockedAnalyserService
			}(),
			context:            getMockContextWithValidRequestBody(),
			expectedStatusCode: 500,
			expectedResponse:   `"someErr"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewPageAnalyserHandler(tt.analyserService)
			resWriter := utils.HackResponseWriter(tt.context)

			handler.Analyse(tt.context)

			assert.Equal(t, tt.expectedStatusCode, resWriter.Status())
			assert.Equal(t, utils.FormatJsonResponse(tt.expectedResponse), utils.FormatJsonResponse(resWriter.Body.String()))
		})
	}
}

func getMockContextWithValidRequestBody() *gin.Context {
	ctx := getMockContext()
	request := model.PageAnalyseRequest{PageURL: "test.com"}
	bytesReq, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Request, _ = http.NewRequest("POST", "/mockURL", bytes.NewReader(bytesReq))
	return ctx
}

func getMockContext() *gin.Context {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest("POST", "/mockURL", nil)
	return c
}
