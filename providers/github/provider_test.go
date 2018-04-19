package github

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/github/model"

	. "gopkg.in/check.v1"
	rcore "gopkg.in/src-d/core-retrieval.v0"
)

func Test(t *testing.T) {
	TestingT(t)
}

type GithubProviderSuite struct {
	DB       *sql.DB
	provider core.RepoProvider
}

var _ = Suite(&GithubProviderSuite{})

func (s *GithubProviderSuite) SetUpTest(c *C) {
	DB := rcore.Database()
	s.DB = DB

	err := core.DropTables(DB, core.GithubProviderName)
	c.Assert(err, IsNil)
	err = core.CreateGithubTable(DB)
	c.Assert(err, IsNil)

	s.provider = NewProvider(core.Config.Github.Token, s.DB)

}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart(c *C) {
	for i := 0; i < 101; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		c.Assert(len(repoUrl.Aliases), Equals, 3)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}
}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart_Repos(c *C) {
	for i := 0; i < 100; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}

	rs, err := model.NewRepositoryStore(s.DB).Find(model.NewRepositoryQuery())
	c.Assert(err, IsNil)
	repos, err := rs.All()
	c.Assert(err, IsNil)

	c.Assert(len(repos), Equals, 100)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_FromStart_ReposTwoPages(c *C) {
	for i := 0; i < 101; i++ {
		repoUrl, err := s.provider.Next()
		c.Assert(err, IsNil)
		c.Assert(repoUrl, NotNil)
		err = s.provider.Ack(nil)
		c.Assert(err, IsNil)
	}

	rs, err := model.NewRepositoryStore(s.DB).Find(model.NewRepositoryQuery())
	c.Assert(err, IsNil)
	repos, err := rs.All()
	c.Assert(err, IsNil)

	c.Assert(len(repos), Equals, 200)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_End(c *C) {
	repo := model.NewRepository()
	repo.GithubID = 999999999999

	repos := []*model.Repository{
		repo,
	}

	// Simulate Ack
	githubProvider, ok := s.provider.(*provider)
	c.Assert(ok, Equals, true)
	err := githubProvider.saveRepos(repos)
	c.Assert(err, IsNil)

	repoUrl, err := s.provider.Next()
	c.Assert(repoUrl, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *GithubProviderSuite) TestGithubProvider_Next_Retry(c *C) {
	repoUrl, err := s.provider.Next()
	c.Assert(err, IsNil)
	c.Assert(repoUrl, NotNil)

	// Simulate an error
	s.provider.Ack(errors.New("WOOPS"))
	repoUrl2, err := s.provider.Next()

	c.Assert(err, IsNil)
	c.Assert(repoUrl, NotNil)
	c.Assert(repoUrl, DeepEquals, repoUrl2)
}

func (s *GithubProviderSuite) TestGithubProviderNilRepositories(c *C) {
	const response string = `
[
    null,
    {
	"id": 40908251,
	"name": "rovers",
	"full_name": "src-d/rovers",
	"owner": {
	  "login": "src-d",
	  "id": 15128793,
	  "avatar_url": "https://avatars2.githubusercontent.com/u/15128793?v=4",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/src-d",
	  "html_url": "https://github.com/src-d",
	  "followers_url": "https://api.github.com/users/src-d/followers",
	  "following_url": "https://api.github.com/users/src-d/following{/other_user}",
	  "gists_url": "https://api.github.com/users/src-d/gists{/gist_id}",
	  "starred_url": "https://api.github.com/users/src-d/starred{/owner}{/repo}",
	  "subscriptions_url": "https://api.github.com/users/src-d/subscriptions",
	  "organizations_url": "https://api.github.com/users/src-d/orgs",
	  "repos_url": "https://api.github.com/users/src-d/repos",
	  "events_url": "https://api.github.com/users/src-d/events{/privacy}",
	  "received_events_url": "https://api.github.com/users/src-d/received_events",
	  "type": "Organization",
	  "site_admin": false
	},
	"private": false,
	"html_url": "https://github.com/src-d/rovers",
	"description": "Rovers is a service to retrieve repository URLs from multiple repository hosting providers.",
	"fork": false,
	"url": "https://api.github.com/repos/src-d/rovers",
	"forks_url": "https://api.github.com/repos/src-d/rovers/forks",
	"keys_url": "https://api.github.com/repos/src-d/rovers/keys{/key_id}",
	"collaborators_url": "https://api.github.com/repos/src-d/rovers/collaborators{/collaborator}",
	"teams_url": "https://api.github.com/repos/src-d/rovers/teams",
	"hooks_url": "https://api.github.com/repos/src-d/rovers/hooks",
	"issue_events_url": "https://api.github.com/repos/src-d/rovers/issues/events{/number}",
	"events_url": "https://api.github.com/repos/src-d/rovers/events",
	"assignees_url": "https://api.github.com/repos/src-d/rovers/assignees{/user}",
	"branches_url": "https://api.github.com/repos/src-d/rovers/branches{/branch}",
	"tags_url": "https://api.github.com/repos/src-d/rovers/tags",
	"blobs_url": "https://api.github.com/repos/src-d/rovers/git/blobs{/sha}",
	"git_tags_url": "https://api.github.com/repos/src-d/rovers/git/tags{/sha}",
	"git_refs_url": "https://api.github.com/repos/src-d/rovers/git/refs{/sha}",
	"trees_url": "https://api.github.com/repos/src-d/rovers/git/trees{/sha}",
	"statuses_url": "https://api.github.com/repos/src-d/rovers/statuses/{sha}",
	"languages_url": "https://api.github.com/repos/src-d/rovers/languages",
	"stargazers_url": "https://api.github.com/repos/src-d/rovers/stargazers",
	"contributors_url": "https://api.github.com/repos/src-d/rovers/contributors",
	"subscribers_url": "https://api.github.com/repos/src-d/rovers/subscribers",
	"subscription_url": "https://api.github.com/repos/src-d/rovers/subscription",
	"commits_url": "https://api.github.com/repos/src-d/rovers/commits{/sha}",
	"git_commits_url": "https://api.github.com/repos/src-d/rovers/git/commits{/sha}",
	"comments_url": "https://api.github.com/repos/src-d/rovers/comments{/number}",
	"issue_comment_url": "https://api.github.com/repos/src-d/rovers/issues/comments{/number}",
	"contents_url": "https://api.github.com/repos/src-d/rovers/contents/{+path}",
	"compare_url": "https://api.github.com/repos/src-d/rovers/compare/{base}...{head}",
	"merges_url": "https://api.github.com/repos/src-d/rovers/merges",
	"archive_url": "https://api.github.com/repos/src-d/rovers/{archive_format}{/ref}",
	"downloads_url": "https://api.github.com/repos/src-d/rovers/downloads",
	"issues_url": "https://api.github.com/repos/src-d/rovers/issues{/number}",
	"pulls_url": "https://api.github.com/repos/src-d/rovers/pulls{/number}",
	"milestones_url": "https://api.github.com/repos/src-d/rovers/milestones{/number}",
	"notifications_url": "https://api.github.com/repos/src-d/rovers/notifications{?since,all,participating}",
	"labels_url": "https://api.github.com/repos/src-d/rovers/labels{/name}",
	"releases_url": "https://api.github.com/repos/src-d/rovers/releases{/id}",
	"deployments_url": "https://api.github.com/repos/src-d/rovers/deployments",
	"created_at": "2015-08-17T15:59:59Z",
	"updated_at": "2018-03-22T08:44:49Z",
	"pushed_at": "2018-04-19T09:11:13Z",
	"git_url": "git://github.com/src-d/rovers.git",
	"ssh_url": "git@github.com:src-d/rovers.git",
	"clone_url": "https://github.com/src-d/rovers.git",
	"svn_url": "https://github.com/src-d/rovers",
	"homepage": "",
	"size": 4171,
	"stargazers_count": 4,
	"watchers_count": 4,
	"language": "HTML",
	"has_issues": true,
	"has_projects": false,
	"has_downloads": true,
	"has_wiki": true,
	"has_pages": false,
	"forks_count": 9,
	"mirror_url": null,
	"archived": false,
	"open_issues_count": 8,
	"license": {
	  "key": "gpl-3.0",
	  "name": "GNU General Public License v3.0",
	  "spdx_id": "GPL-3.0",
	  "url": "https://api.github.com/licenses/gpl-3.0"
	},
	"forks": 9,
	"open_issues": 8,
	"watchers": 4,
	"default_branch": "master",
	"organization": {
	  "login": "src-d",
	  "id": 15128793,
	  "avatar_url": "https://avatars2.githubusercontent.com/u/15128793?v=4",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/src-d",
	  "html_url": "https://github.com/src-d",
	  "followers_url": "https://api.github.com/users/src-d/followers",
	  "following_url": "https://api.github.com/users/src-d/following{/other_user}",
	  "gists_url": "https://api.github.com/users/src-d/gists{/gist_id}",
	  "starred_url": "https://api.github.com/users/src-d/starred{/owner}{/repo}",
	  "subscriptions_url": "https://api.github.com/users/src-d/subscriptions",
	  "organizations_url": "https://api.github.com/users/src-d/orgs",
	  "repos_url": "https://api.github.com/users/src-d/repos",
	  "events_url": "https://api.github.com/users/src-d/events{/privacy}",
	  "received_events_url": "https://api.github.com/users/src-d/received_events",
	  "type": "Organization",
	  "site_admin": false
	},
	"network_count": 9,
	"subscribers_count": 15
      },
      null
]
`

	mockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(rateLimitLimitHeader, "123456789")
		w.Header().Set(rateLimitRemainingHeader, "0")
		w.Write([]byte(response))
	}))
	defer mockedServer.Close()

	s.provider.(*provider).apiClient.endpoint = mockedServer.URL + "?%d=foo"
	defer func() {
		s.provider.(*provider).apiClient.endpoint = githubApiURL
	}()

	mention, err := s.provider.Next()
	c.Assert(err, IsNil)
	c.Assert(mention, NotNil)
	c.Assert(mention.Endpoint, Equals, "git://github.com/src-d/rovers")
	err = s.provider.Ack(nil)
	c.Assert(err, IsNil)
}
