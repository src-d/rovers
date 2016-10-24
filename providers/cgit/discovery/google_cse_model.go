package discovery

type result struct {
	Kind string `json:"kind"`
	URL  struct {
		Type     string `json:"type"`
		Template string `json:"template"`
	} `json:"url"`
	Queries struct {
		Request []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
		} `json:"request"`
		NextPage []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
		} `json:"nextPage"`
	} `json:"queries"`
	Context struct {
		Title string `json:"title"`
	} `json:"context"`
	SearchInformation struct {
		SearchTime            float64 `json:"searchTime"`
		FormattedSearchTime   string  `json:"formattedSearchTime"`
		TotalResults          string  `json:"totalResults"`
		FormattedTotalResults string  `json:"formattedTotalResults"`
	} `json:"searchInformation"`
	Items []item `json:"items"`
}

type item struct {
	Kind             string `json:"kind"`
	Title            string `json:"title"`
	HTMLTitle        string `json:"htmlTitle"`
	Link             string `json:"link"`
	DisplayLink      string `json:"displayLink"`
	Snippet          string `json:"snippet"`
	HTMLSnippet      string `json:"htmlSnippet"`
	CacheID          string `json:"cacheId"`
	FormattedURL     string `json:"formattedUrl"`
	HTMLFormattedURL string `json:"htmlFormattedUrl"`
}

type badRequest struct {
	Error struct {
		Errors []struct {
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
			Message string `json:"message"`
		} `json:"errors"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
