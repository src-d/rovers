package commands

import (
	"time"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-domain/models"
	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/client"
	"github.com/tyba/srcd-rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
)

var Expired = (30 * 24 * time.Hour).Seconds()

// CmdAugur fetches info from Augur for all emails on sourced.people.
//
// NOTE: Augur rate limit is 1 req/s.
type CmdAugur struct {
	CmdBase

	client       *readers.AugurInsightsAPI
	personStore  *models.PersonStore
	emailStore   *social.AugurEmailStore
	insightStore *social.AugurInsightStore
	emails       map[string]time.Time
	now          time.Time
}

func (cmd *CmdAugur) Execute(args []string) error {
	cmd.ChangeLogLevel()

	cmd.client = readers.NewAugurInsightsAPI(client.NewClient(false))
	cmd.personStore = container.GetDomainModelsPersonStore()
	cmd.emailStore = container.GetDomainModelsSocialAugurEmailStore()
	cmd.insightStore = container.GetDomainModelsSocialAugurInsightStore()
	cmd.emails = make(map[string]time.Time)

	start := time.Now()
	defer log15.Info("Done", "elapsed", time.Since(start))

	return cmd.run()
}

func (cmd *CmdAugur) run() error {
	var steps = []struct {
		Name string
		Fn   func() error
	}{
		{"populateEmailSet", cmd.populateEmailSet},
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
	q := cmd.emailStore.Query()
	set, err := cmd.emailStore.Find(q)
	if err != nil {
		return err
	}

	defer log15.Info("Set populated", "total_emails", len(cmd.emails))
	return set.ForEach(func(email *social.AugurEmail) error {
		cmd.emails[email.Email] = email.Last
		return nil
	})
}

func (cmd *CmdAugur) isUpToDate(email string) bool {
	last, ok := cmd.emails[email]
	return ok && time.Since(last).Seconds() < Expired
}

func (cmd *CmdAugur) fetchAugurData() error {
	cmd.now = time.Now()
	q := cmd.personStore.Query()
	set, err := cmd.personStore.Find(q)
	if err != nil {
		return err
	}

	return set.ForEach(func(person *models.Person) error {
		for _, email := range person.Email {
			if cmd.isUpToDate(email) {
				log15.Info("Already up to date",
					"email", email,
					"last_update", cmd.emails[email],
				)
				continue
			}

			if err := cmd.processEmail(email); err != nil {
				log15.Error("Process email failed",
					"email", email,
					"error", err,
				)
				continue
			}
		}
		return nil
	})
}

func (cmd *CmdAugur) processEmail(email string) error {
	insight, resp, err := cmd.client.SearchByEmail(email)
	if err != nil && err != readers.ErrPartialResponse {
		return err
	}

	if err := cmd.setStatus(email, resp.StatusCode); err != nil {
		return err
	}
	return cmd.saveAugurInsights(insight)
}

func (cmd *CmdAugur) setStatus(email string, status int) error {
	doc := cmd.emailStore.New()
	doc.Email = email
	doc.Status = status
	doc.Last = time.Now()
	cmd.emails[email] = doc.Last

	log15.Debug("setStatus", "email", doc.Email, "status", status)
	_, err := cmd.emailStore.Save(doc)
	return err
}

func (cmd *CmdAugur) saveAugurInsights(doc *social.AugurInsight) error {
	return cmd.insightStore.Insert(doc)
}
