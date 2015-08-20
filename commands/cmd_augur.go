package commands

import (
	"time"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-domain/models/rovers/augur"
	"github.com/tyba/srcd-rovers/client"
	"github.com/tyba/srcd-rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
	op "gopkg.in/tyba/storable.v1/operators"
)

// Due to Augur having a rate limit of 1 req/s this is a single goroutine
// process.
type CmdAugur struct {
	FilterBy int `short:"f" long:"filter" description:"filter by status"`
	// SortBy   string `short:"s" long:"sort" default:"email" description:"order by"`
	// Source   string `short:"" long:"source" default:"people" description:""`

	client       *readers.AugurInsightsAPI
	emailStore   *augur.EmailStore
	insightStore *augur.InsightStore
}

func (cmd *CmdAugur) Execute(args []string) error {
	// switch cmd.Source {
	// case "people":
	// 	cmd.emailSource = augur.NewAugurPeopleSource()
	// }

	cmd.client = readers.NewAugurInsightsAPI(client.NewClient(false))
	cmd.emailStore = container.GetDomainModelsRoversAugurEmailStore()
	cmd.insightStore = container.GetDomainModelsRoversAugurInsightStore()

	cmd.process()

	return nil
}

func (cmd *CmdAugur) process() {
	q := cmd.emailStore.Query()
	q.FindWithoutStatus()
	if cmd.FilterBy != 0 {
		q.AddCriteria(op.Eq(augur.Schema.Email.Status, cmd.FilterBy))
	}

	set := cmd.emailStore.MustFind(q)
	defer set.Close()

	for set.Next() {
		email, err := set.Get()
		if err != nil {
			log15.Error("ResultSet.Get", "error", err)
		}
		if err := cmd.processEmail(email); err != nil {
			log15.Error("processEmail", "error", err)
		}
	}
}

func (cmd *CmdAugur) processEmail(e *augur.Email) error {
	insight, resp, err := cmd.client.SearchByEmail(e.Email)
	if err != nil && resp == nil {
		return err
	}

	cmd.setStatus(e, resp.StatusCode)

	if resp.StatusCode == 200 {
		cmd.saveAugurInsights(insight)
		return nil
	}

	return err
}

func (cmd *CmdAugur) setStatus(doc *augur.Email, status int) error {
	doc.Status = status
	doc.Last = time.Now()

	log15.Info("Done", "email", doc.Email, "status", status)
	_, err := cmd.emailStore.Save(doc)
	return err
}

func (cmd *CmdAugur) saveAugurInsights(doc *augur.Insight) error {
	return cmd.insightStore.Insert(doc)
}
