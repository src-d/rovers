package bitbucket

type response struct {
	Pagelen      int          `json:"pagelen"`
	Repositories repositories `json:"values"`
	Next         string       `json:"next"`
}

type repositories []*bitbucketRepository

type bitbucketRepository struct {
	Next    string
	Scm     string `json:"scm" bson:",omitempty"`
	Website string `json:"website" bson:",omitempty"`
	HasWiki bool   `json:"has_wiki" bson:"-"`
	Name    string `json:"name" bson:",omitempty"`
	Links   struct {
		Watchers struct {
			Href string `json:"href" bson:"-"`
		} `json:"watchers" bson:"-"`
		Branches struct {
			Href string `json:"href" bson:"-"`
		} `json:"branches" bson:"-"`
		Tags struct {
			Href string `json:"href" bson:"-"`
		} `json:"tags" bson:"-"`
		Commits struct {
			Href string `json:"href" bson:"-"`
		} `json:"commits" bson:"-"`
		Clone []struct {
			Href string `json:"href" bson:",omitempty"`
			Name string `json:"name" bson:",omitempty"`
		} `json:"clone" bson:",omitempty"`
		Self struct {
			Href string `json:"href" bson:"-"`
		} `json:"self" bson:"-"`
		HTML struct {
			Href string `json:"href" bson:"-"`
		} `json:"html" bson:"-"`
		Avatar struct {
			Href string `json:"href" bson:"-"`
		} `json:"avatar" bson:"-"`
		Hooks struct {
			Href string `json:"href" bson:"-"`
		} `json:"hooks" bson:"-"`
		Forks struct {
			Href string `json:"href" bson:"-"`
		} `json:"forks" bson:"-"`
		Downloads struct {
			Href string `json:"href" bson:"-"`
		} `json:"downloads" bson:"-"`
		Pullrequests struct {
			Href string `json:"href" bson:"-"`
		} `json:"pullrequests" bson:"-"`
	} `json:"links" bson:",omitempty"`
	ForkPolicy string `json:"fork_policy" bson:",omitempty"`
	UUID       string `json:"uuid" bson:",omitempty"`
	Language   string `json:"language" bson:",omitempty"`
	CreatedOn  string `json:"created_on" bson:",omitempty"`
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
	FullName  string `json:"full_name" bson:",omitempty"`
	HasIssues bool   `json:"has_issues" bson:",omitempty"`
	Owner     struct {
		Username    string `json:"username" bson:",omitempty"`
		DisplayName string `json:"display_name" bson:",omitempty"`
		Type        string `json:"type" bson:",omitempty"`
		UUID        string `json:"uuid" bson:",omitempty"`
		Links       struct {
			Self struct {
				Href string `json:"href" bson:"-"`
			} `json:"self"  bson:"-"`
			HTML struct {
				Href string `json:"href"  bson:"-"`
			} `json:"html"  bson:"-"`
			Avatar struct {
				Href string `json:"href"  bson:"-"`
			} `json:"avatar"  bson:"-"`
		} `json:"links"  bson:"-"`
	} `json:"owner" bson:",omitempty"`
	UpdatedOn   string `json:"updated_on" bson:",omitempty"`
	Size        int    `json:"size" bson:",omitempty"`
	Type        string `json:"type" bson:",omitempty"`
	Slug        string `json:"slug" bson:",omitempty"`
	IsPrivate   bool   `json:"is_private" bson:",omitempty"`
	Description string `json:"description" bson:",omitempty"`
}
