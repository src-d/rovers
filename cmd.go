package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tyba/opensource-search/sources/social/commands"
)

func main() {
	parser := flags.NewNamedParser("crawler", flags.Default)
	parser.AddCommand("augur", "Augur retriever", "", &commands.CmdAugur{})
	parser.AddCommand("linkedin", "LinkedIn crawler", "", &commands.CmdLinkedIn{})
	parser.AddCommand("github", "Github crawler", "", &commands.CmdGithub{})
	parser.AddCommand("github-api", "Github API crawler", "", &commands.CmdGithubApi{})
	parser.AddCommand("twitter", "Twitter crawler", "", &commands.CmdTwitter{})
	parser.AddCommand("bitbucket", "Bitbucket API retriever", "", &commands.CmdBitbucket{})
	parser.AddCommand("playground", "Diferent crawler tests", "", &commands.CmdPlayground{})

	_, err := parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}
