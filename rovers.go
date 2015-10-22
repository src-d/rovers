package main

import (
	"os"

	"github.com/src-d/rovers/commands"

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
	_, err = parser.AddCommand("github-profiles", "Github web crawler", "", &commands.CmdGithub{})
	PanicIf(err)
	_, err = parser.AddCommand("github-repos", "Github API repository crawler", "", &commands.CmdGithubApi{})
	PanicIf(err)
	_, err = parser.AddCommand("github-users", "Github API users crawler", "", &commands.CmdGithubApiUsers{})
	PanicIf(err)
	_, err = parser.AddCommand("twitter", "Twitter web crawler", "", &commands.CmdTwitter{})
	PanicIf(err)
	_, err = parser.AddCommand("bitbucket", "Bitbucket API repository crawler", "", &commands.CmdBitbucket{})
	PanicIf(err)
	cmd, err := parser.AddCommand("linkedin", "LinkedIn Company Employees crawler", "", &commands.CmdLinkedIn{})
	PanicIf(err)
	_, err = cmd.AddCommand("no-employees", "Ofelia job: Run 'linkedin' command for just added companies", "", &commands.CmdLinkedInNoEmployees{})
	PanicIf(err)

	_, err = parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}
