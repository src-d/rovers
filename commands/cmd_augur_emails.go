package commands

import (
	"errors"

	"github.com/tyba/srcd-domain/container"
	"github.com/tyba/srcd-domain/models/social"
	"github.com/tyba/srcd-rovers/readers"

	"gopkg.in/inconshreveable/log15.v2"
)

type CmdAugurEmails struct {
	Source string `short:"" long:"source" default:"people" description:""`
	File   string `short:"" long:"file" default:"" description:"requires source=file, path to file"`

	emailSource readers.AugurEmailSource
	emailStore  *social.AugurEmailStore
}

func (cmd *CmdAugurEmails) Execute(args []string) error {
	switch cmd.Source {
	case "people":
		cmd.emailSource = readers.NewAugurPeopleSource()
	case "file":
		if cmd.File != "" {
			return errors.New("no file param provided")
		}
		cmd.emailSource = readers.NewAugurFileSource(cmd.File)
	}

	cmd.emailStore = container.GetDomainModelsSocialAugurEmailStore()

	cmd.populateEmails()

	return nil
}

func (cmd *CmdAugurEmails) populateEmails() {
	for cmd.emailSource.Next() {
		email, err := cmd.emailSource.Get()
		if err != nil {
			log15.Error("ResultSet.Get failed", "error", err)
			continue
		}

		if cmd.isInserted(email) {
			log15.Info("Already inserted", "email", email)
			continue
		}

		if err := cmd.insertEmail(email); err != nil {
			log15.Error("Insert failed", "email", email, "error", err)
		}
	}
}

func (cmd *CmdAugurEmails) isInserted(email string) bool {
	q := cmd.emailStore.Query()
	q.FindByEmail(email)
	doc, err := cmd.emailStore.FindOne(q)
	if err != nil {
		return false
	}
	return doc != nil
}

func (c *CmdAugurEmails) insertEmail(email string) error {
	doc := cmd.emailStore.New()
	doc.Email = email

	return cmd.emailStore.Insert(doc)
}
