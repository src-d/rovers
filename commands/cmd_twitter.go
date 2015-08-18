package commands

import (
	"fmt"

	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/http"
	"github.com/tyba/srcd-rovers/readers"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type augurData struct {
	Profiles struct {
		LinkedInURL string `bson:"linkedin_url"`
		GithubURL   string `bson:"github_url"`
		TwitterURL  string `bson:"twitter_url"`
	}
}

type CmdTwitter struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`

	twitter *readers.TwitterReader
	augur   *mgo.Collection
	storage *mgo.Collection
}

func (t *CmdTwitter) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + t.MongoDBHost)

	t.twitter = readers.NewTwitterReader(http.NewCachedClient(session))
	t.storage = session.DB("sources").C("twitter")
	t.augur = session.DB("sources").C("augur")

	pending := t.get()
	for {
		result := &augurData{}
		if !pending.Next(result) {
			break
		}

		t.processData(result)
	}

	return nil
}

func (t *CmdTwitter) get() *mgo.Iter {
	q := bson.M{
		"profiles.twitter_url": bson.M{
			"$exists": 1,
		},
		"crawler.twitter_url": bson.M{
			"$exists": 0,
		},
	}

	return t.augur.Find(q).Sort("-_id").Iter()
}

func (t *CmdTwitter) processData(d *augurData) {
	url := d.Profiles.TwitterURL
	if t.has(url) {
		fmt.Printf("SKIP: %q\n", url)
		t.done(url, 200)

		return
	}

	p, err := t.twitter.GetProfileByURL(url)
	if err != nil {
		fmt.Printf("ERROR: %q, %s\n", url, err)
		t.done(url, 500)

		return
	}

	t.saveTwitterProfile(p)
	fmt.Printf("DONE: %s\n", p.FullName)
	t.done(url, 200)

	return
}

func (t *CmdTwitter) has(url string) bool {
	q := bson.M{"url": url}

	if c, _ := t.storage.Find(q).Count(); c == 0 {
		return false
	}

	return true
}

func (t *CmdTwitter) done(url string, status int) {
	q := bson.M{"profiles.twitter_url": url}
	s := bson.M{
		"$set": bson.M{
			"crawler.twitter_url": 200,
		},
	}

	_, err := t.augur.UpdateAll(q, s)
	if err != nil {
		panic(err)
	}
}

func (t *CmdTwitter) saveTwitterProfile(p *social.TwitterProfile) error {
	return t.storage.Insert(p)
}
