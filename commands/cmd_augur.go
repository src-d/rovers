package commands

import (
	"time"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/client"
	"github.com/tyba/srcd-rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
)

// Due to Augur having a rate limit of 1 req/s this is a single goroutine
// process.
type CmdAugur struct {
	FilterBy int    `short:"f" long:"filter" description:"filter by status"`
	Source   string `short:"" long:"source" default:"people" description:""`

	emailSource  readers.AugurEmailSource
	client       *readers.AugurInsightsAPI
	emailStore   *social.AugurEmailStore
	insightStore *social.AugurInsightStore
}

func (cmd *CmdAugur) Execute(args []string) error {
	switch cmd.Source {
	case "people":
		cmd.emailSource = augur.NewAugurPeopleSource()
	case "file":
		cmd.emailSource = augur.NewAugurFileSource()
	}

	cmd.client = readers.NewAugurInsightsAPI(client.NewClient(false))
	cmd.emailStore = container.GetDomainModelsSocialAugurEmailStore()
	cmd.insightStore = container.GetDomainModelsSocialAugurInsightStore()

	cmd.process()

	return nil
}

func (cmd *CmdAugur) process() {
	for cmd.emailSource.Next() {
		email, err := cmd.emailSource.Get()
		if err != nil {
			log15.Error("ResultSet.Get", "error", err)
			continue
		}
		if err := cmd.processEmail(email); err != nil {
			log15.Error("processEmail", "error", err)
		}
	}
}

func (cmd *CmdAugur) processEmail(e *social.AugurEmail) error {
	insight, resp, err := cmd.client.SearchByEmail(e.Email)
	if err != nil && resp == nil {
		return err
	}

	if err := cmd.setStatus(e, resp.StatusCode); err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return err
	}

	if err := cmd.saveAugurInsights(insight); err != nil {
		return err
	}

	return nil
}

func (cmd *CmdAugur) setStatus(doc *social.AugurEmail, status int) error {
	doc.Status = status
	doc.Last = time.Now()

	log15.Info("Done", "email", doc.Email, "status", status)
	_, err := cmd.emailStore.Save(doc)
	return err
}

func (cmd *CmdAugur) saveAugurInsights(doc *social.AugurInsight) error {
	return cmd.insightStore.Insert(doc)
}
