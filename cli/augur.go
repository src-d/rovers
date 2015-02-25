package cli

import (
	"fmt"
	"sync"
	"time"

	"github.com/tyba/opensource-search/sources/social/http"
	"github.com/tyba/opensource-search/sources/social/sources"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Augur struct {
	FilterBy    int    `short:"f" long:"filter" description:"filter by status"`
	SortBy      string `short:"s" long:"sort" default:"email" description:"order by"`
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`
	MaxThreads  int    `short:"t" long:"threads" default:"4" description:"number of t"`

	augur        *sources.Augur
	collection   *mgo.Collection
	storage      *mgo.Collection
	emailChannel emailChannel
	sync.WaitGroup
	sync.Mutex
}

type email struct {
	Email  string `bson:"_id"`
	Status int
}

type emailChannel chan *email

func (a *Augur) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + a.MongoDBHost)

	a.augur = sources.NewAugur(http.NewClient(false))
	a.collection = session.DB("social").C("emails")
	a.storage = session.DB("social").C("augur")
	a.emailChannel = make(emailChannel, a.MaxThreads)

	go a.queue()
	a.process()

	return nil
}

func (a *Augur) queue() {
	pending := a.get()
	defer pending.Close()

	for {
		result := &email{}
		if !pending.Next(result) {
			break
		}

		a.emailChannel <- result
	}

	close(a.emailChannel)
}

func (a *Augur) get() *mgo.Iter {
	q := bson.M{
		"status": bson.M{
			"$exists": 0,
		},
	}

	if a.FilterBy != 0 {
		q["status"] = a.FilterBy
	}

	return a.collection.Find(q).Sort(a.SortBy).Iter()
}

func (a *Augur) process() {
	for i := 0; i < a.MaxThreads; i++ {
		a.Add(1)
		go func(i int) {
			defer a.Done()
			for {
				if !a.readFromChannel(i) {
					break
				}
			}
		}(i)
	}

	a.Wait()
}

func (a *Augur) readFromChannel(i int) bool {
	a.Lock()
	email, ok := <-a.emailChannel
	a.Unlock()

	if ok {
		if err := a.processEmail(email); err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}

	return ok
}

func (a *Augur) processEmail(e *email) error {
	r, res, err := a.augur.SearchByEmail(e.Email)
	if err != nil && res == nil {
		return err
	}

	a.setStatus(e, res.StatusCode)

	if res.StatusCode == 200 {
		a.saveAugurInsights(r)
		return nil
	}

	return err
}

func (a *Augur) setStatus(e *email, status int) error {
	q := bson.M{"_id": e.Email}
	s := bson.M{
		"$set": bson.M{
			"status": status,
			"last":   time.Now(),
		},
	}

	fmt.Printf("DONE: %s, %d\n", e.Email, status)

	return a.collection.Update(q, s)
}

func (a *Augur) saveAugurInsights(i *sources.AugurInsights) error {
	return a.storage.Insert(i)
}
