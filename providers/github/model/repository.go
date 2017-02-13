package model

import "github.com/src-d/go-kallax"

type Repository struct {
	kallax.Model      `json:"-" table:"github"`
	kallax.Timestamps `json:"-" kallax:",inline"`

	GithubID    int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Owner       Owner  `json:"owner"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	Fork        bool   `json:"fork"`

	// API URLs
	URL              string `json:"url" kallax:"-"`
	ForksURL         string `json:"forks_url" kallax:"-"`
	KeysURL          string `json:"keys_url" kallax:"-"`
	CollaboratorsURL string `json:"collaborators_url" kallax:"-"`
	TeamsURL         string `json:"teams_url" kallax:"-"`
	HooksURL         string `json:"hooks_url" kallax:"-"`
	IssueEventsURL   string `json:"issue_events_url" kallax:"-"`
	EventsURL        string `json:"events_url" kallax:"-"`
	AssigneesURL     string `json:"assignees_url" kallax:"-"`
	BranchesURL      string `json:"branches_url" kallax:"-"`
	TagsURL          string `json:"tags_url" kallax:"-"`
	BlobsURL         string `json:"blobs_url" kallax:"-"`
	GitTagsURL       string `json:"git_tags_url" kallax:"-"`
	GitRefsURL       string `json:"git_refs_url" kallax:"-"`
	TreesURL         string `json:"trees_url" kallax:"-"`
	StatusesURL      string `json:"statuses_url" kallax:"-"`
	LanguagesURL     string `json:"languages_url" kallax:"-"`
	StargazersURL    string `json:"stargazers_url" kallax:"-"`
	ContributorsURL  string `json:"contributors_url" kallax:"-"`
	SubscribersURL   string `json:"subscribers_url" kallax:"-"`
	SubscriptionURL  string `json:"subscription_url" kallax:"-"`
	CommitsURL       string `json:"commits_url" kallax:"-"`
	GitCommitsURL    string `json:"git_commits_url" kallax:"-"`
	CommentsURL      string `json:"comments_url" kallax:"-"`
	IssueCommentURL  string `json:"issue_comment_url" kallax:"-"`
	ContentsURL      string `json:"contents_url" kallax:"-"`
	CompareURL       string `json:"compare_url" kallax:"-"`
	MergesURL        string `json:"merges_url" kallax:"-"`
	ArchiveURL       string `json:"archive_url" kallax:"-"`
	DownloadsURL     string `json:"downloads_url" kallax:"-"`
	IssuesURL        string `json:"issues_url" kallax:"-"`
	PullsURL         string `json:"pulls_url" kallax:"-"`
	MilestonesURL    string `json:"milestones_url" kallax:"-"`
	NotificationsURL string `json:"notifications_url" kallax:"-"`
	LabelsURL        string `json:"labels_url" kallax:"-"`
	ReleasesURL      string `json:"releases_url" kallax:"-"`
	DeploymentsURL   string `json:"deployments_url" kallax:"-"`
}
