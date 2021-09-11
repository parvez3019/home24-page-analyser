package model

// PageAnalysisResponse represents the model for page analysis response
type PageAnalysisResponse struct {
	HTMLVersion          string         `json:"html_version"`
	Title                string         `json:"title"`
	HeaderCount          map[string]int `json:"header_count"`
	Links                LinksResponse  `json:"links"`
	HasLoginForm         bool           `json:"has_login_form"`
	HasPasswordInputType bool           `json:"-"`
	HasSubmitTypeInput   bool           `json:"-"`
}

// LinkCountResponse represents the model for link/url count response
type LinkCountResponse struct {
	Count int      `json:"count"`
	URLs  []string `json:"urls"`
}

// LinksResponse represents the model for aggregate links response
type LinksResponse struct {
	InternalLinks     LinkCountResponse `json:"internal"`
	ExternalLinks     LinkCountResponse `json:"external"`
	InaccessibleLinks LinkCountResponse `json:"inaccessible"`
}
