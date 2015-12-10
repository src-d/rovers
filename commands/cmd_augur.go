package commands

import (
	"time"

	"github.com/src-d/domain/container"
	"github.com/src-d/domain/models"
	"github.com/src-d/domain/models/social"
	"github.com/src-d/rovers/client"
	"github.com/src-d/rovers/metrics"
	"github.com/src-d/rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/src-d/storable.v1"
)

var Expired = (30 * 24 * time.Hour).Seconds()

// CmdAugur fetches up to 1_000_000 user insights a month of a set of users of
// our people database by their email.
type CmdAugur struct {
	CmdBase

	client       *readers.AugurInsightsAPI
	personStore  *models.PersonStore
	insightStore *social.AugurInsightStore
	emailSet     map[string]bool // Set of all emails + Done status
}

func (c *CmdAugur) Execute(args []string) error {
	start := time.Now()

	c.ChangeLogLevel()
	c.client = readers.NewAugurInsightsAPI(client.NewClient(true))
	c.personStore = container.GetDomainModelsPersonStore()
	c.insightStore = container.GetDomainModelsSocialAugurInsightStore()
	c.emailSet = make(map[string]bool)

	err := c.run()
	log15.Info("Done", "elapsed", time.Since(start))

	return err
}

// run performs all steps required to update our Augur collection:
//
// 1. Start with a fresh empty set, name: `email-set`.
//
// 2. Gather all emails from `sourced.people` and insert them into `email-set`.
//
// 3. Prune `email-set` by removing all emails from `sources.augur` with `done=true`.
//
// 4. Fetch Augur data for all remaining emails in `email-set`.
func (c *CmdAugur) run() error {
	var steps = []struct {
		Name string
		Fn   func() error
	}{
		{"populateEmailSet", c.populateEmailSet},
		{"pruneEmailSet", c.pruneEmailSet},
		{"fetchAugurData", c.fetchAugurData},
	}
	for _, step := range steps {
		start := time.Now()
		err := step.Fn()
		elapsed := time.Since(start)
		if err != nil {
			log15.Error("Done", "step", step.Name, "elapsed", elapsed)
			return err
		} else {
			log15.Info("Done", "step", step.Name, "elapsed", elapsed)
		}
	}
	return nil
}

func (c *CmdAugur) populateEmailSet() error {
	q := c.personStore.Query()
	set, err := c.personStore.Find(q)
	if err != nil {
		return err
	}

	err = set.ForEach(func(person *models.Person) error {
		for _, email := range person.Email {
			c.emailSet[email.Address] = false
		}
		return nil
	})

	log15.Info("Set populated", "total_emails", len(c.emailSet))
	return err
}

func (c *CmdAugur) pruneEmailSet() error {
	q := c.insightStore.Query()
	q.FindDone()
	set, err := c.insightStore.Find(q)
	if err != nil {
		return err
	}

	deleted := 0
	err = set.ForEach(func(insight *social.AugurInsight) error {
		if insight.Done {
			deleted++
			delete(c.emailSet, insight.InputEmail)
		}
		return nil
	})

	log15.Info("Set pruned", "pruned_emails", deleted)
	return nil
}

func (c *CmdAugur) fetchAugurData() error {
	for email, done := range c.emailSet {
		if done {
			continue
		}
		c.emailSet[email] = true

		err := c.processEmail(email)
		if err == readers.ErrRateLimitExceeded {
			log15.Warn("Rate limit reached. Stopping...")
			return storable.ErrStop
		}
		if err != nil {
			log15.Error("Process email failed",
				"email", email,
				"error", err,
			)
			continue
		}
		metrics.AugurProcessed.Inc()
	}
	return nil
}

func (c *CmdAugur) processEmail(email string) error {
	insight, resp, err := c.client.SearchByEmail(email)
	if err != nil && err != readers.ErrPartialResponse {
		return err
	}
	if resp == nil {
		return err
	}
	if insight == nil {
		log15.Debug("Empty insight",
			"email", email,
			"status_code", resp.StatusCode,
			"error", err,
		)
		return nil
	}
	return c.saveAugurInsights(insight)
}

func (c *CmdAugur) saveAugurInsights(doc *social.AugurInsight) error {
	return c.insightStore.Insert(doc)
}
