package commands

import (
	"time"

	"github.com/src-d/domain/container"
	"github.com/src-d/domain/models/social"
	"github.com/src-d/rovers/metrics"
	"github.com/src-d/rovers/readers"

	"github.com/mcuadros/go-github/github"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/src-d/storable.v1"
)

type CmdGitHubAPIRepos struct {
	github  *readers.GithubAPI
	storage *social.GithubRepositoryStore
}

func (c *CmdGitHubAPIRepos) Execute(args []string) error {
	c.github = readers.NewGithubAPI()
	c.storage = container.GetDomainModelsSocialGithubRepositoryStore()

	start := time.Now()
	since := c.getSince()
	for {
		log15.Info("Requesting repositories...", "since", since)

		repos, resp, err := c.github.GetAllRepositories(since)
		if err != nil {
			return err
		}

		c.save(repos)
		if resp.NextPage == 0 && resp.NextPage == since {
			break
		}

		since = resp.NextPage
	}

	log15.Info("Done", "elapsed", time.Since(start))
	return nil
}

func (c *CmdGitHubAPIRepos) getSince() int {
	q := c.storage.Query()

	q.Sort(storable.Sort{{social.Schema.GithubRepository.GithubID, storable.Desc}})
	repo, err := c.storage.FindOne(q)
	if err != nil {
		return 0
	}

	return repo.GithubID
}

func (c *CmdGitHubAPIRepos) getRepositories(since int) (
	repos []github.Repository, resp *github.Response, err error,
) {
	metrics.GitHubReposRequested.Inc()

	start := time.Now()
	repos, resp, err = c.github.GetAllRepositories(since)
	if err != nil {
		log15.Error("GetAllRepositories failed",
			"since", since,
			"error", err,
		)
		metrics.GitHubReposFailed.WithLabelValues("ghapi_request").Inc()
		return
	}

	elapsed := time.Since(start)
	microseconds := float64(elapsed) / float64(time.Microsecond)
	metrics.GitHubReposRequestDur.Observe(microseconds)
	return
}

func (c *CmdGitHubAPIRepos) save(repos []github.Repository) {
	for _, repo := range repos {
		doc := c.createNewDocument(repo)
		if _, err := c.storage.Save(doc); err != nil {
			log15.Error("Repository save failed",
				"repo", doc.FullName,
				"error", err,
			)
			metrics.GitHubReposFailed.WithLabelValues("db_insert").Inc()
		}
	}

	numRepos := len(repos)
	metrics.GitHubReposProcessed.Add(float64(numRepos))
	log15.Info("Repositories saved", "num_repos", numRepos)
}

func (c *CmdGitHubAPIRepos) createNewDocument(repo github.Repository) *social.GithubRepository {
	doc := c.storage.New()
	processGithubRepository(doc, repo)
	return doc
}

func processGithubRepository(doc *social.GithubRepository, repo github.Repository) {
	if repo.ID != nil {
		doc.GithubID = *repo.ID
	}
	if repo.Name != nil {
		doc.Name = *repo.Name
	}
	if repo.FullName != nil {
		doc.FullName = *repo.FullName
	}
	if repo.Description != nil {
		doc.Description = *repo.Description
	}
	if repo.Owner != nil {
		processGithubUser(doc.Owner, *repo.Owner)
		doc.Owner.SetId(bson.NewObjectId())
	}
	if repo.Fork != nil {
		doc.Fork = *repo.Fork
	}
}
