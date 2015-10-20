package commands

import (
	"time"

	"github.com/src-d/domain/container"
	"github.com/src-d/domain/models"
	"github.com/src-d/domain/models/social"
	"github.com/src-d/rovers/client"
	"github.com/src-d/rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/tyba/storable.v1"
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

func (cmd *CmdAugur) Execute(args []string) error {
	start := time.Now()

	cmd.ChangeLogLevel()
	cmd.client = readers.NewAugurInsightsAPI(client.NewClient(true))
	cmd.personStore = container.GetDomainModelsPersonStore()
	cmd.insightStore = container.GetDomainModelsSocialAugurInsightStore()
	cmd.emailSet = make(map[string]bool)

	err := cmd.run()
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
func (cmd *CmdAugur) run() error {
	var steps = []struct {
		Name string
		Fn   func() error
	}{
		{"populateEmailSet", cmd.populateEmailSet},
		{"pruneEmailSet", cmd.pruneEmailSet},
		{"fetchAugurData", cmd.fetchAugurData},
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

func (cmd *CmdAugur) populateEmailSet() error {
	q := cmd.personStore.Query()
	set, err := cmd.personStore.Find(q)
	if err != nil {
		return err
	}

	err = set.ForEach(func(person *models.Person) error {
		for _, email := range person.Email {
			cmd.emailSet[email] = false
		}
		return nil
	})

	log15.Info("Set populated", "total_emails", len(cmd.emailSet))
	return err
}

func (cmd *CmdAugur) pruneEmailSet() error {
	q := cmd.insightStore.Query()
	q.FindDone()
	set, err := cmd.insightStore.Find(q)
	if err != nil {
		return err
	}

	deleted := 0
	err = set.ForEach(func(insight *social.AugurInsight) error {
		if insight.Done {
			deleted++
			delete(cmd.emailSet, insight.InputEmail)
		}
		return nil
	})

	log15.Info("Set pruned", "pruned_emails", deleted)
	return nil
}

func (cmd *CmdAugur) fetchAugurData() error {
	for email, done := range cmd.emailSet {
		if done {
			continue
		}
		cmd.emailSet[email] = true

		// if cmd.isUpToDate(email) {
		// 	log15.Info("Already up to date",
		// 		"email", email,
		// 		"last_update", cmd.emailSet[email],
		// 	)
		// 	continue
		// }

		err := cmd.processEmail(email)
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
	}
	return nil
}

func (cmd *CmdAugur) isUpToDate(email string) bool {
	upToDate, ok := cmd.emailSet[email]
	return ok && upToDate
}

func (cmd *CmdAugur) processEmail(email string) error {
	insight, resp, err := cmd.client.SearchByEmail(email)
	if err != nil && err != readers.ErrPartialResponse {
		return err
	}
	if resp == nil {
		return err
	}
	if insight == nil {
		log15.Info("Empty insight",
			"email", email,
			"status_code", resp.StatusCode,
			"error", err,
		)
		return nil
	}
	return cmd.saveAugurInsights(insight)
}

func (cmd *CmdAugur) saveAugurInsights(doc *social.AugurInsight) error {
	return cmd.insightStore.Insert(doc)
}
