package commands

import (
	"net/url"
	"time"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/client"
	"github.com/tyba/srcd-rovers/readers"

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

func (b *CmdBitbucket) insertRepository(res *readers.BitbucketPagedResult) (n int, err error) {
	for _, value := range res.Values {
		repository := b.store.New()
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

		err = b.store.Insert(repository)
		if err != nil {
			return
		}
		n++

		log15.Debug("Saved repository", "repo", repository)
	}

	log15.Debug("Save", "num_repos", len(res.Values), "next", res.Next)

	return n, nil
}
