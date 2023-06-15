package main

import (
	"fmt"
	"mp/munchies/internal/db"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func makeNutrientsCommand() *cli.Command {
	return &cli.Command{
		Name:    "nutrients",
		Aliases: []string{"n", "nuts"},
		Usage:   "list nutrients",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() > 0 {
				cli.ShowSubcommandHelpAndExit(cCtx, 1)
			}

			return doListNutrients()
		},
	}
}

func doListNutrients() error {
	dB, err := db.New(db.DB_PATH)
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error opening database %q: %s.", db.DB_PATH, err), 1)
	}

	var nuts db.Nutrients
	if err := nuts.ReadFrom(dB); err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error reading nutrients from database %q: %s.", db.DB_PATH, err), 1)

	}

	sort.Slice(nuts, func(i, j int) bool {
		return nuts[i].Name < nuts[j].Name
	})

	for _, nut := range nuts {
		fmt.Println(nut.Name)
	}

	return nil
}
