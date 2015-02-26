package cli

import (
	"fmt"

	"github.com/tyba/opensource-search/sources/social/http"
	"github.com/tyba/opensource-search/sources/social/readers"
	"github.com/tyba/opensource-search/types/social"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Twitter struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`

	twitter *readers.TwitterReader
	augur   *mgo.Collection
	storage *mgo.Collection
}

func (t *Twitter) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + t.MongoDBHost)

	t.twitter = readers.NewTwitterReader(http.NewCachedClient(session))
	t.storage = session.DB("social").C("twitter")
	t.augur = session.DB("social").C("augur")

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

func (t *Twitter) get() *mgo.Iter {
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

func (t *Twitter) processData(d *augurData) {
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

func (t *Twitter) has(url string) bool {
	q := bson.M{"url": url}

	if c, _ := t.storage.Find(q).Count(); c == 0 {
		return false
	}

	return true
}

func (t *Twitter) done(url string, status int) {
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

func (t *Twitter) saveTwitterProfile(p *social.TwitterProfile) error {
	return t.storage.Insert(p)
}
