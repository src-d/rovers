package commands

import (
	"net/url"
	"time"

	"github.com/src-d/rovers/client"
	"github.com/src-d/rovers/metrics"
	"github.com/src-d/rovers/readers"
	"gop.kg/src-d/domain@v3.0/container"
	"gop.kg/src-d/domain@v3.0/models/social"

	"gopkg.in/inconshreveable/log15.v2"
)

type CmdBitbucket struct {
	bitbucket *readers.BitbucketAPI
	store     *social.BitbucketRepositoryStore
}

func (b *CmdBitbucket) Execute(args []string) error {
	b.bitbucket = readers.NewBitbucketAPI(client.NewClient(true))
	b.store = container.GetDomainModelsSocialBitbucketRepositoryStore()

	startExecute := time.Now()
	startThousand := time.Now()

	var (
		next url.Values
		i    = 0
	)
	for {
		result, err := b.bitbucket.GetRepositories(next)
		if err != nil {
			return err
		}

		inserted, err := b.insertRepository(result)
		if err != nil {
			return err
		}

		i = (i + inserted) % 1000
		if i > 1000 {
			i -= 1000
			log15.Info("Saved 1k repositories", "elapsed", time.Since(startThousand))
			startThousand = time.Now()
		}

		next = result.Next.Query()
	}

	log15.Info("Done", "elapsed", time.Since(startExecute))

	return nil
}

func (b *CmdBitbucket) getAfter() string {
	q := b.store.Query()
	repo, err := b.store.FindOne(q)
	if err != nil {
		log15.Error("getAfter query failed")
		metrics.BitbucketFailed.WithLabelValues("getAfter").Inc()
		return ""
	}
	// 2008-11-09T14:59:29.540461+00:00
	return repo.CreatedOn.Format("2006-01-02T15:04:05.999999+07:00")
}

func (b *CmdBitbucket) insertRepository(res *readers.BitbucketPagedResult) (n int, err error) {
	for _, value := range res.Values {
		repository := b.store.New()
		parseRepository(repository, value)

		err = b.store.Insert(repository)
		if err != nil {
			metrics.BitbucketFailed.WithLabelValues("insert").Inc()
			return
		}
		n++

		log15.Debug("Saved repository", "repo", repository)
	}

	metrics.BitbucketProcessed.Inc()
	log15.Debug("Save", "num_repos", len(res.Values), "next", res.Next)

	return n, nil
}

func parseRepository(repository *social.BitbucketRepository, value readers.Repository) {
	repository.CreatedOn = value.CreatedOn
	repository.Description = value.Description
	repository.ForkPolicy = value.ForkPolicy
	repository.FullName = value.FullName
	repository.HasIssues = value.HasIssues
	repository.HasWiki = value.HasWiki
	repository.IsPrivate = value.IsPrivate
	repository.Language = value.Language
	repository.Links.Avatar = value.Links.Avatar.Href
	repository.Links.Clone = value.Links.Clone
	repository.Links.Self = value.Links.Self.Href
	repository.Name = value.Name
	repository.Owner.Links.Avatar = value.Owner.Links.Avatar.Href
	repository.Owner.Links.Html = value.Owner.Links.Html.Href
	repository.Owner.Links.Self = value.Owner.Links.Self.Href
	repository.Owner.DisplayName = value.Owner.DisplayName
	repository.Owner.Type = value.Owner.Type
	repository.Owner.Username = value.Owner.Username
	repository.Owner.UUID = value.Owner.UUID
	repository.Size = value.Size
	repository.Type = value.Type
	repository.UpdatedOn = value.UpdatedOn
	repository.URL = value.Links.Html.Href
	repository.UUID = value.UUID
	repository.VCS = value.SCM
}
