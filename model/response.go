package model

type PageAnalysisResponse struct {
	HTMLVersion          string         `json:"html_version"`
	Title                string         `json:"title"`
	HeaderCount          map[string]int `json:"header_count"`
	Links                LinksResponse  `json:"links"`
	HasLoginForm         bool           `json:"has_login_form"`
	HasPasswordInputType bool           `json:"-"`
	HasSubmitTypeInput   bool           `json:"-"`
}

type LinkCountResponse struct {
	Count int      `json:"count"`
	URLs  []string `json:"urls"`
}

type LinksResponse struct {
	InternalLinks     LinkCountResponse `json:"internal"`
	ExternalLinks     LinkCountResponse `json:"external"`
	InaccessibleLinks LinkCountResponse `json:"inaccessible"`
}
