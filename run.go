package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tyba/opensource-search/sources/social/cli"
)

func main() {
	parser := flags.NewNamedParser("crawler", flags.Default)
	parser.AddCommand("augur", "Augur retriever", "", &cli.Augur{})
	parser.AddCommand("linkedin", "LinkedIn crawler", "", &cli.LinkedIn{})
	parser.AddCommand("github", "Github crawler", "", &cli.Github{})
	parser.AddCommand("twitter", "Twitter crawler", "", &cli.Twitter{})
	parser.AddCommand("bitbucket", "Bitbucket API retriever", "", &cli.Bitbucket{})
	parser.AddCommand("playground", "Diferent crawler tests", "", &cli.Playground{})

	_, err := parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}
