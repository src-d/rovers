package commands

import (
	"fmt"
	"time"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-rovers/metrics"
	"github.com/tyba/srcd-rovers/readers"

	"github.com/mcuadros/go-github/github"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
)

type CmdGithubApiUsers struct {
	github  *readers.GithubAPI
	storage *mgo.Collection
}

func (cmd *CmdGithubApiUsers) Execute(args []string) error {
	defer metrics.Push()

	session := container.GetMgoSession()
	defer session.Close()
	cmd.github = readers.NewGithubAPI()
	cmd.storage = session.DB("github").C("users.api")

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		seconds := float64(elapsed) / float64(time.Microsecond)
		metrics.GitHubUsersTotalDur.Observe(seconds)
		log15.Info("Done", "elapsed", elapsed)
	}()

	since := cmd.getSince()
	for {
		fmt.Printf("Requesting since %d ...", since)
		users, resp, err := cmd.github.GetAllUsers(since)
		if err != nil {
			return err
		}

		cmd.save(users)
		if resp.NextPage == 0 && resp.NextPage == since {
			break
		}

		since = resp.NextPage
	}

	return nil
}

func (cmd *CmdGithubApiUsers) getSince() int {
	var user *github.User
	cmd.storage.Find(nil).Sort("-id").One(&user)

	if user == nil {
		return 0
	}

	return *user.ID
}

func (cmd *CmdGithubApiUsers) save(users []github.User) {
	for _, user := range users {
		if err := cmd.storage.Insert(user); err != nil {
			log15.Error("User save failed",
				"user", *user.Name,
				"error", err,
			)
			labels := []string{"error", err.Error()}
			if user.Name != nil {
				labels = append(labels, []string{"user", *user.Name}...)
			}
			metrics.GitHubUsersFailed.WithLabelValues(labels...).Inc()
		}
	}

	metrics.GitHubUsersProcessed.Add(float64(len(users)))
	log15.Info(fmt.Sprintf("Saved %d users", len(users)))
}
