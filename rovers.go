package main

import (
	"os"

	"github.com/tyba/srcd-rovers/commands"

	"github.com/jessevdk/go-flags"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	parser := flags.NewNamedParser("rovers", flags.Default)

	var err error
	_, err = parser.AddCommand("augur", "Augur Insights API crawler", "", &commands.CmdAugur{})
	PanicIf(err)
	_, err = parser.AddCommand("augur-emails", "Populates social.augur.emails", "", &commands.CmdAugurEmails{})
	PanicIf(err)
	_, err = parser.AddCommand("github", "Github web crawler", "", &commands.CmdGithub{})
	PanicIf(err)
	_, err = parser.AddCommand("github-api", "Github API repository crawler", "", &commands.CmdGithubApi{})
	PanicIf(err)
	_, err = parser.AddCommand("github-api-users", "Github API users crawler", "", &commands.CmdGithubApiUsers{})
	PanicIf(err)
	_, err = parser.AddCommand("twitter", "Twitter web crawler", "", &commands.CmdTwitter{})
	PanicIf(err)
	_, err = parser.AddCommand("bitbucket", "Bitbucket API repository crawler", "", &commands.CmdBitbucket{})
	PanicIf(err)

	_, err = parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}
