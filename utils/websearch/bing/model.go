package bing

type result struct {
	Type     string `json:"_type"`
	WebPages struct {
		WebSearchURL          string  `json:"webSearchUrl"`
		TotalEstimatedMatches int     `json:"totalEstimatedMatches"`
		Values                []value `json:"value"`
	} `json:"webPages"`
	RankingResponse struct {
		Mainline struct {
			Items []struct {
				AnswerType  string `json:"answerType"`
				ResultIndex int    `json:"resultIndex"`
				Value       struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"mainline"`
	} `json:"rankingResponse"`
}

type value struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	DisplayURL      string `json:"displayUrl"`
	Snippet         string `json:"snippet"`
	DateLastCrawled string `json:"dateLastCrawled"`
}
