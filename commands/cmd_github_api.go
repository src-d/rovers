package commands

import (
	"fmt"

	"github.com/mcuadros/go-github/github"
	"github.com/tyba/oss/sources/social/readers"

	"gopkg.in/mgo.v2"
)

type CmdGithubApi struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`

	github  *readers.GithubAPIReader
	storage *mgo.Collection
}

func (l *CmdGithubApi) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + l.MongoDBHost)

	l.github = readers.NewGithubAPIReader(nil)
	l.storage = session.DB("github").C("repositories")

	since := l.getSince()
	for {
		fmt.Printf("Requesting since %d ...", since)
		repos, resp, err := l.github.GetAllRepositories(since)
		if err != nil {
			return err
		}

		l.save(repos)
		if resp.NextPage == 0 && resp.NextPage == since {
			break
		}

		since = resp.NextPage
	}

	return nil
}

func (l *CmdGithubApi) getSince() int {
	var r github.Repository
	l.storage.Find(nil).Sort("-id").One(&r)

	return *r.ID
}

func (l *CmdGithubApi) save(repos []github.Repository) {
	for _, r := range repos {
		if err := l.storage.Insert(r); err != nil {
			fmt.Println("error", err)
		}
	}

	fmt.Printf("saved %d repositories\n", len(repos))
}
