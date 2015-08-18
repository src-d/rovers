package commands

import (
	"fmt"
	"net/url"

	"github.com/tyba/srcd-rovers/http"
	"github.com/tyba/srcd-rovers/readers"

	"gopkg.in/mgo.v2"
)

type CmdBitbucket struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`
	MaxThreads  int    `short:"t" long:"threads" default:"4" description:"number of t"`

	bitbucket *readers.BitbucketReader
	storage   *mgo.Collection
}

func (b *CmdBitbucket) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + b.MongoDBHost)

	b.bitbucket = readers.NewBitbucketReader(http.NewClient(true))
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
	fmt.Printf("Retrieved: %d repositorie(s)\nNext: %s\n", len(res.Values), res.Next)

	for _, r := range res.Values {
		if err := b.storage.Insert(r); err != nil {
			return err
		}
	}

	return nil
}
