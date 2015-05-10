package commands

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tyba/oss/domain/models/social"
	"github.com/tyba/oss/sources/social/http"
	"github.com/tyba/oss/sources/social/readers"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CmdGithub struct {
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`
	MaxThreads  int    `short:"t" long:"threads" default:"4" description:"number of t"`

	github  *readers.GithubReader
	augur   *mgo.Collection
	storage *mgo.Collection

	c chan *githubUrlData
	sync.WaitGroup
	sync.Mutex
}

type githubUrlData struct {
	Url string
}

func (l *CmdGithub) Execute(args []string) error {
	l.c = make(chan *githubUrlData, l.MaxThreads)

	session, _ := mgo.Dial("mongodb://" + l.MongoDBHost)

	l.github = readers.NewGithubReader(http.NewCachedClient(session))
	l.storage = session.DB("github").C("profiles")
	l.augur = session.DB("github").C("urls")

	go l.queue()
	l.process()

	return nil
}

func (l *CmdGithub) queue() {
	pending := l.get()
	defer pending.Close()

	for {
		result := &githubUrlData{}
		if !pending.Next(result) {
			break
		}

		if pending.Err() != nil {
			fmt.Println(pending.Err())
			break
		}

		l.c <- result
	}

	close(l.c)
}

func (l *CmdGithub) get() *mgo.Iter {
	q := bson.M{
		"done": bson.M{
			"$exists": 0,
		},
	}

	return l.augur.Find(q).Iter()
}

func (l *CmdGithub) process() {
	for i := 0; i < l.MaxThreads; i++ {
		l.Add(1)
		go func() {
			defer l.Done()
			for {
				l.Lock()
				url, _ := <-l.c
				l.Unlock()
				l.processData(url)
			}
		}()
	}

	l.Wait()
}

func (l *CmdGithub) processData(d *githubUrlData) {
	if d == nil {
		fmt.Println("Empty")
		return
	}

	url := strings.Replace(d.Url, "https:", "http:", 1)
	if l.has(url) {
		fmt.Printf("SKIP: %q\n", url)
		l.done(url, 200)
		return
	}

	p, err := l.github.GetProfileByURL(url)
	if err != nil {
		if err == http.NotFound {
			l.done(url, 404)
		} else {
			l.done(url, 500)
		}

		fmt.Printf("ERROR: %q, %s\n", url, err)
		return
	}

	if err := l.saveGithubProfile(p); err != nil {
		fmt.Printf("ERROR saving: %q, %s\n", url, err)
		return
	}

	fmt.Printf("DONE: Organization: %b Username: %s\n", p.Organization, p.Username)
	l.done(url, 200)

	return
}

func (l *CmdGithub) has(url string) bool {
	q := bson.M{"url": url}

	if c, _ := l.storage.Find(q).Count(); c == 0 {
		return false
	}

	return true
}

func (l *CmdGithub) done(url string, status int) {
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

func (l *CmdGithub) saveGithubProfile(p *social.GithubProfile) error {
	p.SetId(bson.NewObjectId())
	return l.storage.Insert(p)
}
