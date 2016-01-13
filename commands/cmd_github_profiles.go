package commands

import (
	"fmt"
	"strings"
	"sync"

	"github.com/src-d/rovers/client"
	"github.com/src-d/rovers/readers"
	"gop.kg/src-d/domain@v3.0/container"
	"gop.kg/src-d/domain@v3.0/models/social"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CmdGitHubProfiles struct {
	CmdBase

	MaxThreads int `short:"t" long:"threads" default:"4" description:"number of t"`

	github  *readers.GithubWebCrawler
	augur   *mgo.Collection
	storage *mgo.Collection

	c chan *githubUrlData
	sync.WaitGroup
	sync.Mutex
}

type githubUrlData struct {
	Url string
}

func (c *CmdGitHubProfiles) Execute(args []string) error {
	c.c = make(chan *githubUrlData, c.MaxThreads)

	session := container.GetMgoSession()
	defer session.Close()
	c.github = readers.NewGithubWebCrawler(client.NewClient(true))
	c.storage = session.DB("github").C("profiles")
	c.augur = session.DB("github").C("urls")

	go c.queue()
	c.process()

	return nil
}

func (c *CmdGitHubProfiles) queue() {
	pending := c.get()
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

		c.c <- result
	}

	close(c.c)
}

func (c *CmdGitHubProfiles) get() *mgo.Iter {
	q := bson.M{
		"done": bson.M{"$exists": 0},
	}
	return c.augur.Find(q).Iter()
}

func (c *CmdGitHubProfiles) process() {
	c.Add(c.MaxThreads)
	for i := 0; i < c.MaxThreads; i++ {
		go func() {
			defer c.Done()
			for {
				url, _ := <-c.c
				c.processData(url)
			}
		}()
	}

	c.Wait()
}

func (c *CmdGitHubProfiles) processData(d *githubUrlData) {
	if d == nil {
		fmt.Println("Empty")
		return
	}

	url := strings.Replace(d.Url, "https:", "http:", 1)
	if c.has(url) {
		fmt.Printf("SKIP: %q\n", url)
		c.done(url, 200)
		return
	}

	p, err := c.github.GetProfileByURL(url)
	if err != nil {
		if err == client.NotFound {
			c.done(url, 404)
		} else {
			c.done(url, 500)
		}

		fmt.Printf("ERROR: %q, %s\n", url, err)
		return
	}

	if err := c.saveGithubProfile(p); err != nil {
		fmt.Printf("ERROR saving: %q, %s\n", url, err)
		return
	}

	fmt.Printf("DONE: Organization: %b Username: %s\n", p.Organization, p.Username)
	c.done(url, 200)

	return
}

func (c *CmdGitHubProfiles) has(url string) bool {
	q := bson.M{"url": url}

	if c, _ := c.storage.Find(q).Count(); c == 0 {
		return false
	}

	return true
}

func (c *CmdGitHubProfiles) done(url string, status int) {
	q := bson.M{"url": url}
	s := bson.M{
		"$set": bson.M{
			"done": status,
		},
	}

	_, err := c.augur.UpdateAll(q, s)
	if err != nil {
		panic(err)
	}
}

func (c *CmdGitHubProfiles) saveGithubProfile(p *social.GithubProfile) error {
	p.SetId(bson.NewObjectId())
	return c.storage.Insert(p)
}
