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
	SetupExternalLinkAccessibilityDownstreamMocks()

	_ = client.
		Post("/analyse").
		BodyString(`{"page_url": "https://www.test.com"}`).
		Expect(t).
		Status(200).
		JSON(`{
			"has_login_form":true,
			"header_count":{"h1":1,"h2":2,"h3":3,"h4":4,"h5":5,"h6":6},
			"html_version":"HTML 3",
			"links":{
			"external":{"count":4,"urls":["https://www.facebook.com","https://www.google.com","https://www.amazon.com","https://www.example.com/inaccessible"]},
			"internal":{"count":2,"urls":["/internal1","/internal2"]},
			"inaccessible":{"count":1,"urls":["https://www.example.com/inaccessible"]}},
			"title":"Parvez Hassan Test Page"}`).
		Done()
}

func SetupTestDomainDownstreamMocks() {
	httpmock.RegisterResponder("GET", "https://www.test.com", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, testHTMLPageResponse), nil
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
	<a href="/internal1"></a>
	<a href="/internal2"></a>
	<a href="https://www.facebook.com"></a>
	<a href="https://www.google.com"></a>
	<a href="https://www.amazon.com"></a>
	<a href="https://www.example.com/inaccessible"></a>
	<form class="modal-content animate" action="/action_page.php">
	<div class="container">
	  <label for="uname"><b>Username</b></label>
	  <input type="text" placeholder="Enter Username" name="uname" required>
	
	  <label for="psw"><b>Password</b></label>
	  <input type="password" placeholder="Enter Password" name="psw" required>
	  <button type="submit">Login</button>
	  <label>
		<input type="checkbox" checked="checked" name="remember"> Remember me
	  </label>
	</div>
	</form>
</body>
</html>
`
