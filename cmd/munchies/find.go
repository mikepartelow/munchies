package main

import (
	"fmt"
	"mp/munchies/internal/db"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func makeFindCommand() *cli.Command {
	return &cli.Command{
		Name:        "find",
		Aliases:     []string{"f"},
		Usage:       "find a food",
		ArgsUsage:   "<search string>",
		Description: "find lists nutrient details if there is only a single match, otherwise, it lists the names of all matching foods.",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 1 {
				cli.ShowSubcommandHelpAndExit(cCtx, 1)
			}

			return doFind(cCtx.Args().First())
		},
	}
}

func doFind(term string) error {
	dB, err := getDb()
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error opening database %q: %s.", getDbPath(), err), 1)
	}
	defer dB.Close()

	var foods db.Foods
	if err := foods.Match(term, dB); err != nil {
		if err == db.NoMatch {
			return cli.Exit(fmt.Sprintf("no matches for %q.", term), 1)
		}
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error opening database %q: %s.", getDbPath(), err), 1)
	}

	if len(foods) == 1 {
		fmt.Println(foods[0].Name)
		fmt.Println("---")
		for _, fnut := range foods[0].Nutrients {
			fmt.Println("  " + fnut.Name)
		}
	} else {
		for _, food := range foods {
			fmt.Println(food.Name)
		}
	}

	return nil
}
