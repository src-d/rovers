package commands

import (
	"time"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/metrics"
	"github.com/tyba/srcd-rovers/readers"

	"github.com/mcuadros/go-github/github"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/tyba/storable.v1"
)

type CmdGithubApiUsers struct {
	github  *readers.GithubAPI
	storage *social.GithubUserStore
}

func (c *CmdGithubApiUsers) Execute(args []string) error {
	defer metrics.Push()

	c.github = readers.NewGithubAPI()
	c.storage = container.GetDomainModelsSocialGithubUserStore()

	start := time.Now()
	defer log15.Info("Done", "elapsed", time.Since(start))

	since := c.getSince()
	for {
		log15.Info("Requesting users...", "since", since)

		users, resp, err := c.getUsers(since)
		if err != nil {
			return err
		}
		c.save(users)

		if resp.NextPage == 0 && resp.NextPage == since {
			break
		}

		since = resp.NextPage
	}

	return nil
}

func (c *CmdGithubApiUsers) getSince() int {
	q := c.storage.Query()
	q.Sort(storable.Sort{{social.Schema.GithubUser.GithubID, storable.Desc}})
	user, err := c.storage.FindOne(q)
	if err != nil {
		log15.Error("getSince query failed")
		return 0
	}

	return user.GithubID
}

func (c *CmdGithubApiUsers) getUsers(since int) (
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

func (c *CmdGithubApiUsers) save(users []github.User) {
	for _, user := range users {
		doc := c.createNewDocument(user)
		if _, err := c.storage.Save(doc); err != nil {
			log15.Error("User save failed",
				"user", doc.Login,
				"error", err,
			)
			metrics.GitHubUsersFailed.WithLabelValues("db_insert").Inc()
		}
	}

	numUsers := len(users)
	metrics.GitHubUsersProcessed.Add(float64(numUsers))
	log15.Info("Users saved", "num_users", numUsers)
}

func (c *CmdGithubApiUsers) createNewDocument(user github.User) *social.GithubUser {
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
