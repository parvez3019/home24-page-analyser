package model

type PageAnalysisResponse struct {
	HTMLVersion       string            `json:"html_version"`
	Title             string            `json:"title"`
	HeaderCount       map[string]int    `json:"header_count"`
	Links             LinksResponse     `json:"links"`
	InaccessibleLinks LinkCountResponse `json:"inaccessible_links"`
	HasLoginForm      bool              `json:"has_login_form"`
}

type LinkCountResponse struct {
	Count int      `json:"count"`
	URLs  []string `json:"urls"`
}

type LinksResponse struct {
	InternalLinks LinkCountResponse `json:"internal"`
	ExternalLinks LinkCountResponse `json:"external"`
}
