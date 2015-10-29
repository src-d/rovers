package commands

import (
	"github.com/src-d/domain/container"
	"github.com/src-d/domain/models/social"
	"github.com/src-d/rovers/client"
	"github.com/src-d/rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2/bson"
)

type CmdTwitter struct {
	augur   *social.AugurInsightStore
	store   *social.TwitterProfileStore
	twitter *readers.TwitterReader
}

func (cmd *CmdTwitter) Execute(args []string) error {
	cmd.twitter = readers.NewTwitterReader(client.NewClient(true))
	cmd.store = container.GetDomainModelsSocialTwitterProfileStore()
	cmd.augur = container.GetDomainModelsSocialAugurInsightStore()

	q := cmd.augur.Query()
	q.FindTwitterHandles()

	pending, err := cmd.augur.Find(q)
	if err != nil {
		return err
	}
	for pending.Next() {
		result, err := pending.Get()
		if err != nil {
			break
		}

		cmd.fetchProfiles(result)
	}

	return nil
}

func (cmd *CmdTwitter) fetchProfiles(insight *social.AugurInsight) {
	url := social.NormalizeTwitterURL(insight.Profiles.Handle)
	if cmd.isStored(url) {
		log15.Info("Skipping", "url", url)
		cmd.flagAsDone(insight, 200)
		return
	}

	profile, err := cmd.twitter.GetProfileByURL(url)
	if err != nil {
		log15.Error("No profile found", "url", url, "error", err)
		cmd.flagAsDone(insight, 500)
		return
	}

	err = cmd.store.Insert(profile)
	if err != nil {
		log15.Error("Couldn't insert profile", "url", url, "error", err)
		cmd.flagAsDone(insight, 500)
		return
	}

	log15.Info("Done", "name", profile.FullName)
	cmd.flagAsDone(insight, 200)
}

func (cmd *CmdTwitter) isStored(url string) bool {
	q := cmd.store.Query()
	q.FindByURL(url)

	_, err := cmd.store.FindOne(q)
	return err == nil
}

func (cmd *CmdTwitter) flagAsDone(insight *social.AugurInsight, status int) {
	q := cmd.augur.Query()
	q.FindByTwitterHandle(insight.Profiles.Handle)

	update := bson.M{"done": true}
	cmd.augur.RawUpdate(q, update, true) // el error para Rita
}
