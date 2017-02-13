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
	_, err = parser.AddCommand("repos", "Get repos from internet", "", &commands.CmdRepoProviders{})
	PanicIf(err)
	_, err = parser.AddCommand("initdb", "Create tables",
		"Create the necessary tables used by the providers into the database", &commands.CmdCreateTables{})
	PanicIf(err)

	_, err = parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}

	select {}
}
