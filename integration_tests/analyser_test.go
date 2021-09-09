package integration_tests

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func Test_ShouldReturnExpectedParsedStatisticResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	SetupTestDomainDownstreamMocks()
	_ = client.
		Post("/analyse").
		BodyString(`{"page_url": "https://www.test.com"}`).
		Expect(t).
		Status(200).
		JSON(`{
			"has_login_form":false,
			"header_count":{"h1":1,"h2":2,"h3":3,"h4":4,"h5":5,"h6":6},
			"html_version":"HTML 3",
			"inaccessible_links":{"count":0,"urls":null},
			"links":{"external":{"count":0,"urls":null},
			"internal":{"count":0,"urls":null}},
			"title":""}`).
		Done()
}

func SetupTestDomainDownstreamMocks() {
	httpmock.RegisterResponder("GET", "https://www.test.com", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, testHTMLPageResponse), nil
	})
}

var testHTMLPageResponse = `
<!DOCTYPE HTML PUBLIC "HTML 3">
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
	<title>Parvez Hassan Test Page</title>
</head>
<body>
	<h1>Heading</h1>
	<h2>Heading</h2>
	<h2>Heading</h2>
	<h3>Heading</h3>
	<h3>Heading</h3>
	<h3>Heading</h3>
	<h4>Heading</h4>
	<h4>Heading</h4>
	<h4>Heading</h4>
	<h4>Heading</h4>
	<h5>Heading</h5>
	<h5>Heading</h5>
	<h5>Heading</h5>
	<h5>Heading</h5>
	<h5>Heading</h5>
	<h6>Heading</h6>
	<h6>Heading</h6>
	<h6>Heading</h6>
	<h6>Heading</h6>
	<h6>Heading</h6>
	<h6>Heading</h6>
	<a href="/internal"></a>
	<a href="https://www.example.com/external"></a>
	<a href="https://www.example.com/inaccessible"></a>
</body>
</html>
`
