package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/tyba/opensource-search/sources/social/http"
	"github.com/tyba/opensource-search/sources/social/readers"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CmdAugur struct {
	FilterBy    int    `short:"f" long:"filter" description:"filter by status"`
	SortBy      string `short:"s" long:"sort" default:"email" description:"order by"`
	MongoDBHost string `short:"m" long:"mongo" default:"localhost" description:"mongodb hostname"`
	MaxThreads  int    `short:"t" long:"threads" default:"4" description:"number of t"`

	augur        *readers.AugurReader
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

func (a *CmdAugur) Execute(args []string) error {
	session, _ := mgo.Dial("mongodb://" + a.MongoDBHost)

	a.augur = readers.NewAugurReader(http.NewClient(false))
	a.collection = session.DB("sources").C("emails")
	a.storage = session.DB("sources").C("augur")
	a.emailChannel = make(emailChannel, a.MaxThreads)

	go a.queue()
	a.process()

	return nil
}

func (a *CmdAugur) queue() {
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

func (a *CmdAugur) get() *mgo.Iter {
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

func (a *CmdAugur) process() {
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

func (a *CmdAugur) readFromChannel(i int) bool {
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

func (a *CmdAugur) processEmail(e *email) error {
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

func (a *CmdAugur) setStatus(e *email, status int) error {
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

func (a *CmdAugur) saveAugurInsights(i *readers.AugurInsights) error {
	return a.storage.Insert(i)
}
