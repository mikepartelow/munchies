package main

import (
	"fmt"
	"mp/munchies/internal/db"
	"mp/munchies/pkg/usda"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)

	app := &cli.App{
		Name:  "munchies",
		Usage: "learn stuff about food",
		Commands: []*cli.Command{
			makeFindCommand(),
			makeIngestCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}

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

func makeIngestCommand() *cli.Command {
	return &cli.Command{
		Name:        "ingest",
		Aliases:     []string{"i"},
		Usage:       "ingest USDA JSONs",
		ArgsUsage:   "<path to USDA JSONs> <path to database>",
		Description: "injest creates a new database <path to database> and imports the foods found at <path to USDA JSON>.",
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() < 2 {
				cli.ShowSubcommandHelpAndExit(cCtx, 1)
			}

			return doIngest(cCtx.Args().Get(0), cCtx.Args().Get(1))
		},
	}
}

func doFind(term string) error {
	return nil
}

func doIngest(pathToUSDAJson, pathToDatabase string) error {
	foods, err := usda.MustRead(pathToUSDAJson)
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error reading USDA JSON files at %q: %s.", pathToUSDAJson, err), 1)
	}
	_ = os.RemoveAll(pathToDatabase)
	dB, err := db.New(pathToDatabase)
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error creating database %q: %s.", pathToDatabase, err), 1)
	}
	if err := dB.Migrate(); err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error migrating database %q: %s.", pathToDatabase, err), 1)
	}

	defer dB.Close()

	units := make(map[string]struct{})
	nutrients := make(map[string]struct{})

	for _, food := range foods {
		for _, fnut := range food.FoodNutrients {
			uname := fnut.Nutrient.UnitName
			if _, ok := units[uname]; !ok {
				if err := (db.Unit{
					Name: uname,
				}.WriteTo(dB)); err != nil {
					log.Error().Err(err).Send()
					return cli.Exit(fmt.Sprintf("Error writing unit %q to database %q: %s.", uname, pathToUSDAJson, err), 1)
				}
				units[uname] = struct{}{}
			}

			// DRY: so error prone!
			nname := fnut.Nutrient.Name

			if _, ok := nutrients[nname]; !ok {
				if err := (db.Nutrient{
					Name: nname,
				}.WriteTo(dB)); err != nil {
					log.Error().Err(err).Send()
					return cli.Exit(fmt.Sprintf("Error writing unit %q to database %q: %s.", nname, pathToUSDAJson, err), 1)
				}
				nutrients[nname] = struct{}{}
			}
		}
	}

	return nil
}
