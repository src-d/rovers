package commands

import (
	"fmt"

	"github.com/tyba/srcd-rovers/readers"

	"github.com/mcuadros/go-github/github"
	"gopkg.in/mgo.v2"
)

type CmdGithubApiUsers struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`

	github  *readers.GithubAPI
	storage *mgo.Collection
}

func (l *CmdGithubApiUsers) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + l.MongoDBHost)

	l.github = readers.NewGithubAPI()
	l.storage = session.DB("github").C("users.api")

	since := l.getSince()
	for {
		fmt.Printf("Requesting since %d ...", since)
		users, resp, err := l.github.GetAllUsers(since)
		if err != nil {
			return err
		}

		l.save(users)
		if resp.NextPage == 0 && resp.NextPage == since {
			break
		}

		since = resp.NextPage
	}

	return nil
}

func (l *CmdGithubApiUsers) getSince() int {
	var r *github.User
	l.storage.Find(nil).Sort("-id").One(&r)

	if r == nil {
		return 0
	}

	return *r.ID
}

func (l *CmdGithubApiUsers) save(users []github.User) {
	for _, u := range users {
		if err := l.storage.Insert(u); err != nil {
			fmt.Println("error", err)
		}
	}

	fmt.Printf("saved %d repositories\n", len(users))
}
