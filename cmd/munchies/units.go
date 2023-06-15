package main

import (
	"fmt"
	"mp/munchies/internal/db"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func makeUnitsCommand() *cli.Command {
	return &cli.Command{
		Name:    "units",
		Aliases: []string{"u"},
		Usage:   "list units",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 0 {
				cli.ShowSubcommandHelpAndExit(cCtx, 1)
			}

			return doListUnits()
		},
	}
}

func doListUnits() error {
	dB, err := getDb()
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error opening database %q: %s.", getDbPath(), err), 1)
	}

	var units db.Units
	if err := units.ReadFrom(dB); err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error reading units from database %q: %s.", getDbPath(), err), 1)

	}

	sort.Slice(units, func(i, j int) bool {
		return units[i].Name < units[j].Name
	})

	for _, u := range units {
		fmt.Println(u.Name)
	}

	return nil
}
