package commands

import (
	"errors"
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
	File     string `short:"" long:"file" default:"" description:"requires source=file, path to file"`

	emailSource  readers.AugurEmailSource
	client       *readers.AugurInsightsAPI
	emailStore   *social.AugurEmailStore
	insightStore *social.AugurInsightStore
}

func (cmd *CmdAugur) Execute(args []string) error {
	switch cmd.Source {
	case "people":
		cmd.emailSource = readers.NewAugurPeopleSource()
	case "file":
		if cmd.File != "" {
			return errors.New("no file param provided")
		}
		cmd.emailSource = readers.NewAugurFileSource(cmd.File)
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

func (cmd *CmdAugur) processEmail(email string) error {
	insight, resp, err := cmd.client.SearchByEmail(email)
	if err != nil && resp == nil {
		return err
	}

	if err := cmd.setStatus(email, resp.StatusCode); err != nil {
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

func (cmd *CmdAugur) setStatus(email string, status int) error {
	doc := cmd.emailStore.New()
	doc.Status = status
	doc.Last = time.Now()

	log15.Info("Done", "email", doc.Email, "status", status)
	_, err := cmd.emailStore.Save(doc)
	return err
}

func (cmd *CmdAugur) saveAugurInsights(doc *social.AugurInsight) error {
	return cmd.insightStore.Insert(doc)
}
