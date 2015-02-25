package cli

import (
	"fmt"
	"strings"

	"github.com/tyba/opensource-search/sources/social/http"
	"github.com/tyba/opensource-search/sources/social/sources"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Github struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`

	github  *sources.Github
	augur   *mgo.Collection
	storage *mgo.Collection
}

type githubUrlData struct {
	Url string
}

func (l *Github) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + l.MongoDBHost)

	l.github = sources.NewGithub(http.NewCachedClient(session))
	l.storage = session.DB("social").C("github")
	l.augur = session.DB("social").C("github_url")

	pending := l.get()
	for {
		result := &githubUrlData{}
		if !pending.Next(result) {
			break
		}

		l.processData(result)
	}

	return nil
}

func (l *Github) get() *mgo.Iter {
	q := bson.M{
		"done": bson.M{
			"$exists": 1,
		},
	}

	return l.augur.Find(q).Skip(300000).Iter()
}

func (l *Github) processData(d *githubUrlData) {
	url := strings.Replace(d.Url, "https:", "http:", 1)
	if l.has(url) {
		fmt.Printf("SKIP: %q\n", url)
		l.done(url, 200)
		return
	}

	p, err := l.github.GetProfileByURL(url)
	if err != nil {
		fmt.Printf("ERROR: %q, %s\n", url, err)
		l.done(url, 500)
		return
	}

	l.saveGithubProfile(p)
	fmt.Printf("DONE: %s\n", p.Description)
	l.done(url, 200)

	return
}

func (l *Github) has(url string) bool {
	q := bson.M{"url": url}

	if c, _ := l.storage.Find(q).Count(); c == 0 {
		return false
	}

	return true
}

func (l *Github) done(url string, status int) {
	q := bson.M{"url": url}
	s := bson.M{
		"$set": bson.M{
			"done": status,
		},
	}

	_, err := l.augur.UpdateAll(q, s)
	if err != nil {
		panic(err)
	}
}

func (l *Github) saveGithubProfile(p *sources.GithubProfile) error {
	return l.storage.Insert(p)
}
