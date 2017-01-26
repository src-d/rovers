package bitbucket

import "github.com/src-d/go-kallax"

type response struct {
	Pagelen      int          `json:"pagelen"`
	Repositories repositories `json:"values"`
	Next         string       `json:"next"`
}

type repositories []*bitbucketRepository

type bitbucketRepository struct {
	kallax.Model `table:"bitbucket"`
	kallax.Timestamps

	Next    string
	Scm     string `json:"scm"`
	Website string `json:"website"`
	HasWiki bool   `json:"has_wiki" kallax:"-"`
	Name    string `json:"name"`
	Links   struct {
		Watchers struct {
			Href string `json:"href" kallax:"-"`
		} `json:"watchers" kallax:"-"`
		Branches struct {
			Href string `json:"href" kallax:"-"`
		} `json:"branches" kallax:"-"`
		Tags struct {
			Href string `json:"href" kallax:"-"`
		} `json:"tags" kallax:"-"`
		Commits struct {
			Href string `json:"href" kallax:"-"`
		} `json:"commits" kallax:"-"`
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
		Self struct {
			Href string `json:"href" kallax:"-"`
		} `json:"self" kallax:"-"`
		HTML struct {
			Href string `json:"href" kallax:"-"`
		} `json:"html" kallax:"-"`
		Avatar struct {
			Href string `json:"href" kallax:"-"`
		} `json:"avatar" kallax:"-"`
		Hooks struct {
			Href string `json:"href" kallax:"-"`
		} `json:"hooks" kallax:"-"`
		Forks struct {
			Href string `json:"href" kallax:"-"`
		} `json:"forks" kallax:"-"`
		Downloads struct {
			Href string `json:"href" kallax:"-"`
		} `json:"downloads" kallax:"-"`
		Pullrequests struct {
			Href string `json:"href" kallax:"-"`
		} `json:"pullrequests" kallax:"-"`
	} `json:"links"`
	ForkPolicy string `json:"fork_policy"`
	UUID       string `json:"uuid"`
	Language   string `json:"language"`
	CreatedOn  string `json:"created_on"`
	Parent     *struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		UUID     string `json:"uuid"`
	} `json:"parent"`
	FullName  string `json:"full_name"`
	HasIssues bool   `json:"has_issues"`
	Owner     struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Type        string `json:"type"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self struct {
				Href string `json:"href" kallax:"-"`
			} `json:"self"  kallax:"-"`
			HTML struct {
				Href string `json:"href"  kallax:"-"`
			} `json:"html"  kallax:"-"`
			Avatar struct {
				Href string `json:"href"  kallax:"-"`
			} `json:"avatar"  kallax:"-"`
		} `json:"links"  kallax:"-"`
	} `json:"owner"`
	UpdatedOn   string `json:"updated_on"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	Slug        string `json:"slug"`
	IsPrivate   bool   `json:"is_private"`
	Description string `json:"description"`
}
