package service

import (
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"home24-page-analyser/http"
	httpMocks "home24-page-analyser/http/mocks"
	"home24-page-analyser/model"
	"home24-page-analyser/service/html_parser"
	"home24-page-analyser/service/html_parser/mocks"
	"testing"
)

func Test_analyserService_Analyse(t *testing.T) {
	type fields struct {
		client http.Client
		parser html_parser.Parser
	}
	tests := []struct {
		name            string
		fields          fields
		request         model.PageAnalyseRequest
		want            model.PageAnalysisResponse
		wantErr         bool
		wantErrResponse string
	}{
		{
			name: "should return valid page analysis response",
			fields: fields{
				client: setupHTTPDownstreamMocks(),
				parser: setupParserMocks(),
			},
			request: model.PageAnalyseRequest{PageURL: "/someURL"},
			want: model.PageAnalysisResponse{
				Links: model.LinksResponse{
					ExternalLinks:     model.LinkCountResponse{Count: 1, URLs: []string{"/external"}},
					InaccessibleLinks: model.LinkCountResponse{Count: 1, URLs: []string{"/external"}},
				},
			},
			wantErr: false,
		},
		{
			name: "should return error when fails while fetching requested page",
			fields: fields{
				client: func() http.Client {
					mockHttp := &httpMocks.Client{}
					mockHttp.On("Get", "/someURL").
						Return(nil, errors.New("someErr"))
					return mockHttp
				}(),
			},
			request:         model.PageAnalyseRequest{PageURL: "/someURL"},
			want:            model.PageAnalysisResponse{},
			wantErr:         true,
			wantErrResponse: "Error SOMETHING_WENT_WRONG, someErr",
		},
		{
			name: "should return error when fails while parsing the page content",
			fields: fields{
				client: setupHTTPDownstreamMocks(),
				parser: func() html_parser.Parser {
					parser := &mocks.Parser{}
					parser.On("Parse", mock.Anything, "/someURL").
						Return(model.PageAnalysisResponse{}, errors.New("someParsingErr"))
					return parser
				}(),
			},
			request:         model.PageAnalyseRequest{PageURL: "/someURL"},
			want:            model.PageAnalysisResponse{},
			wantErr:         true,
			wantErrResponse: "Error SOMETHING_WENT_WRONG, someParsingErr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnalyzerService(tt.fields.client, tt.fields.parser)
			got, err := a.Analyse(tt.request)
			if tt.wantErr {
				assert.EqualError(t, err, tt.wantErrResponse)
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Analyse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupParserMocks() html_parser.Parser {
	return func() html_parser.Parser {
		parser := &mocks.Parser{}
		parser.On("Parse", mock.Anything, "/someURL").
			Return(model.PageAnalysisResponse{
				Links: model.LinksResponse{
					ExternalLinks:     model.LinkCountResponse{Count: 1, URLs: []string{"/external"}},
					InaccessibleLinks: model.LinkCountResponse{Count: 0, URLs: []string{}},
				},
			}, nil)
		return parser
	}()
}

func setupHTTPDownstreamMocks() http.Client {
	return func() http.Client {
		mockHttp := &httpMocks.Client{}
		mockHttp.On("Get", "/someURL").
			Return(httpmock.NewJsonResponse(200, nil))
		mockHttp.On("Get", "/external").
			Return(httpmock.NewJsonResponse(500, nil))
		return mockHttp
	}()
}
