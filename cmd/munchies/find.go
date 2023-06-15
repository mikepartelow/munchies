package main

import "github.com/urfave/cli/v2"

func makeFindCommand() *cli.Command {
	return &cli.Command{
		Name:        "find",
		Aliases:     []string{"f"},
		Usage:       "find a food",
		ArgsUsage:   "<search string>",
		Description: "find lists nutrient details if there is only a single match, otherwise, it lists the names of all matching foods.",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() < 1 {
				cli.ShowSubcommandHelpAndExit(cCtx, 1)
			}

			return doFind(cCtx.Args().First())
		},
	}
}

func doFind(term string) error {
	return nil
}
