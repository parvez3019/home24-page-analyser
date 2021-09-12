package html_parser

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"home24-page-analyser/model"
	"home24-page-analyser/utils"
	"io"
	"os"
	"testing"
)

func Test_parser_Parse(t *testing.T) {
	type args struct {
		reader io.Reader
		domain string
	}
	tests := []struct {
		name    string
		args    args
		want    model.PageAnalysisResponse
		wantErr bool
	}{
		{
			name: "Should parse the reader content and return valid page anaylsis response",
			args: args{
				reader: readTestHTMLFile(),
				domain: "https://www.test.com",
			},
			want:    expectedValidPageAnalysisResponse(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, err := p.Parse(tt.args.reader, tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func expectedValidPageAnalysisResponse() model.PageAnalysisResponse {
	return model.PageAnalysisResponse{
		HTMLVersion: "HTML 3",
		Title:       "Parvez Hassan Test Page",
		HeaderCount: map[string]int{"h1": 1, "h2": 2, "h3": 3, "h4": 4, "h5": 5, "h6": 6},
		Links: model.LinksResponse{
			InternalLinks: model.LinkCountResponse{
				Count: 2,
				URLs:  []string{"/internal1", "/internal2"},
			},
			ExternalLinks: model.LinkCountResponse{
				Count: 4,
				URLs:  []string{"https://www.facebook.com", "https://www.google.com", "https://www.amazon.com", "https://www.example.com/inaccessible"},
			},
			InaccessibleLinks: model.LinkCountResponse{Count: 0, URLs: []string{}},
		},
		HasLoginForm:         true,
		HasPasswordInputType: true,
		HasSubmitTypeInput:   true,
	}
}

func readTestHTMLFile() io.Reader {
	relativePath := utils.GetRelativePath("/home24-page-analyser/service/html_parser/test_data/test_html_page.html")
	file, err := os.Open(relativePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
