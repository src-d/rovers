package commands

import (
	"net/url"

	"github.com/tyba/srcd-rovers/client"
	"github.com/tyba/srcd-rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
)

type CmdBitbucket struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`

	bitbucket *readers.BitbucketAPI
	storage   *mgo.Collection
}

func (b *CmdBitbucket) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + b.MongoDBHost)

	b.bitbucket = readers.NewBitbucketAPI(client.NewClient(true))
	b.storage = session.DB("sources").C("bitbucket_repositories")

	r, err := b.bitbucket.GetRepositories(url.Values{})
	if err != nil {
		return err
	}

	for {
		r, err = b.bitbucket.GetRepositories(r.Next.Query())
		if err != nil {
			return err
		}

		b.saveBitbucketPagedResult(r)
	}

	return nil
}

func (b *CmdBitbucket) saveBitbucketPagedResult(res *readers.BitbucketPagedResult) error {
	log15.Info("Save", "num_repos", len(res.Values), "next", res.Next)

	for _, r := range res.Values {
		if err := b.storage.Insert(r); err != nil {
			return err
		}
	}

	return nil
}
