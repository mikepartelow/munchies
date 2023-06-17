package main

import (
	"fmt"
	"mp/munchies/pkg/ingest"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func makeIngestCommand() *cli.Command {
	return &cli.Command{
		Name:        "ingest",
		Aliases:     []string{"i"},
		Usage:       "ingest USDA JSONs",
		ArgsUsage:   "<path to USDA JSONs>",
		Description: "injest creates a new database and imports the foods found at <path to USDA JSON>.",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() != 1 {
				cli.ShowSubcommandHelpAndExit(cCtx, 1)
			}

			return doIngest(cCtx.Args().Get(0))
		},
	}
}

func doIngest(usdaJsonPath string) error {
	i, err := ingest.New(usdaJsonPath, getDbPath())
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error ingesting USDA JSON files at %q: %s.", usdaJsonPath, err), 1)
	}
	defer i.DB.Close()

	fmt.Printf("%d units\n%d nutrients\n%d foods.\n", i.Units, i.Nutrients, i.Foods)

	return nil
}
