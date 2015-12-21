package commands

import (
	"time"

	"github.com/src-d/rovers/metrics"
	"github.com/src-d/rovers/readers"
	"gop.kg/src-d/domain@v2.1/container"
	"gop.kg/src-d/domain@v2.1/models/social"

	"github.com/mcuadros/go-github/github"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/src-d/storable.v1"
)

type CmdGitHubAPIUsers struct {
	CmdBase

	github  *readers.GithubAPI
	storage *social.GithubUserStore
}

func (c *CmdGitHubAPIUsers) Execute(args []string) error {
	c.CmdBase.ChangeLogLevel()

	defer metrics.Push()

	c.github = readers.NewGithubAPI()
	c.storage = container.GetDomainModelsSocialGithubUserStore()

	start := time.Now()
	since := c.getSince()
	for {
		log15.Info("Requesting users...", "since", since)

		users, resp, err := c.getUsers(since)
		if err != nil {
			return err
		}
		log15.Debug("Got users...",
			"count", len(users),
			"status", resp.StatusCode,
			"next_page", resp.NextPage,
		)

		if len(users) == 0 {
			log15.Info("No more users. Stopping crawl...")
			break
		}

		c.save(users)

		if resp.NextPage == 0 && resp.NextPage == since {
			break
		}

		since = resp.NextPage
	}

	log15.Info("Done", "elapsed", time.Since(start))
	return nil
}

func (c *CmdGitHubAPIUsers) getSince() int {
	q := c.storage.Query()
	q.Sort(storable.Sort{{social.Schema.GithubUser.GithubID, storable.Desc}})
	user, err := c.storage.FindOne(q)
	if err != nil {
		log15.Crit("getSince query failed")
		return 0
	}

	return user.GithubID
}

func (c *CmdGitHubAPIUsers) getUsers(since int) (
	users []github.User, resp *github.Response, err error,
) {
	metrics.GitHubUsersRequested.Inc()

	start := time.Now()
	users, resp, err = c.github.GetAllUsers(since)
	if err != nil {
		log15.Error("GetAllUsers failed",
			"since", since,
			"error", err,
		)
		metrics.GitHubUsersFailed.WithLabelValues("ghapi_request").Inc()
		return
	}

	elapsed := time.Since(start)
	microseconds := float64(elapsed) / float64(time.Microsecond)
	metrics.GitHubUsersRequestDur.Observe(microseconds)
	return
}

func (c *CmdGitHubAPIUsers) save(users []github.User) {
	saved := 0
	for _, user := range users {
		doc := c.createNewDocument(user)
		if _, err := c.storage.Save(doc); err != nil {
			log15.Error("User save failed",
				"user", doc.Login,
				"error", err,
			)
			metrics.GitHubUsersFailed.WithLabelValues("db_insert").Inc()
			continue
		}
		saved++
	}

	numUsers := len(users)
	metrics.GitHubUsersProcessed.Add(float64(numUsers))
	log15.Info("Users saved",
		"num_input", numUsers,
		"num_saved", saved,
	)
}

func (c *CmdGitHubAPIUsers) createNewDocument(user github.User) *social.GithubUser {
	doc := c.storage.New()
	processGithubUser(doc, user)
	return doc
}

func processGithubUser(doc *social.GithubUser, user github.User) {
	if user.ID != nil {
		doc.GithubID = *user.ID
	}
	if user.Login != nil {
		doc.Login = *user.Login
	}
	if user.AvatarURL != nil {
		doc.Avatar = *user.AvatarURL
	}
	if user.Type != nil {
		doc.Type = *user.Type
	}
}
