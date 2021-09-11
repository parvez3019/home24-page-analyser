package integration_tests

import (
	"github.com/jarcoal/httpmock"
	"home24-page-analyser/utils"
	"net/http"
	"testing"
)

func Test_ShouldReturnExpectedParsedStatisticResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	SetupTestDomainDownstreamMocks()
	SetupExternalLinkAccessibilityDownstreamMocks()

	_ = client.
		Post("/analyse").
		BodyString(`{"page_url": "https://www.test.com"}`).
		Expect(t).
		Status(200).
		JSON(utils.LoadFileAsStringFromPath("/home24-page-analyser/integration_tests/test_data/analyse_response.json")).
		Done()
}

func SetupTestDomainDownstreamMocks() {
	httpmock.RegisterResponder("GET", "https://www.test.com", func(req *http.Request) (*http.Response, error) {
		htmlPage := utils.LoadFileAsStringFromPath("/home24-page-analyser/integration_tests/test_data/test_html_page.html")
		return httpmock.NewStringResponse(200, htmlPage), nil
	})
}

func SetupExternalLinkAccessibilityDownstreamMocks() {
	httpmock.RegisterResponder("GET", "https://www.facebook.com", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, nil)
	})
	httpmock.RegisterResponder("GET", "https://www.google.com", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, nil)
	})
	httpmock.RegisterResponder("GET", "https://www.amazon.com", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, nil)
	})
}
