package commands

import (
	"github.com/src-d/rovers/client"
	"github.com/src-d/rovers/metrics"
	"github.com/src-d/rovers/readers"
	"gop.kg/src-d/domain@v3/container"
	"gop.kg/src-d/domain@v3/models/social"

	"gopkg.in/inconshreveable/log15.v2"
)

type CmdTwitter struct {
	twitter *readers.TwitterReader
	augur   *social.AugurInsightStore
	store   *social.TwitterProfileStore
}

func (cmd *CmdTwitter) Execute(args []string) error {
	cmd.twitter = readers.NewTwitterReader(client.NewClient(true))
	cmd.augur = container.GetDomainModelsSocialAugurInsightStore()
	cmd.store = container.GetDomainModelsSocialTwitterProfileStore()

	q := cmd.augur.Query()
	q.FindTwitterHandles()

	pending, err := cmd.augur.Find(q)
	if err != nil {
		return err
	}

	return pending.ForEach(func(insight *social.AugurInsight) error {
		cmd.fetchProfiles(insight)
		return nil
	})
}

func (cmd *CmdTwitter) fetchProfiles(insight *social.AugurInsight) {
	defer metrics.TwitterProcessed.Inc()

	for _, URL := range insight.TwitterURL {
		if cmd.isStored(URL) {
			log15.Info("Skipping", "url", URL)
			cmd.flagAsDone(insight, 200)
			return
		}

		profile, err := cmd.twitter.GetProfileByURL(URL)
		if err != nil {
			log15.Error("No profile found", "url", URL, "error", err)
			cmd.flagAsDone(insight, 500)
			return
		}

		err = cmd.store.Insert(profile)
		if err != nil {
			log15.Error("Couldn't insert profile", "url", URL, "error", err)
			cmd.flagAsDone(insight, 500)
			return
		}

		log15.Info("Done", "name", profile.FullName)
		cmd.flagAsDone(insight, 200)
	}
}

func (cmd *CmdTwitter) isStored(URL string) bool {
	q := cmd.store.Query()
	q.FindByURL(URL)

	_, err := cmd.store.FindOne(q)
	return err == nil
}

func (cmd *CmdTwitter) flagAsDone(insight *social.AugurInsight, status int) {
	insight.TwitterDone = true
	insight.TwitterStatus = status
	cmd.augur.Save(insight)
}
