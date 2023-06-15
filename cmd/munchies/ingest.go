package main

import (
	"fmt"
	"mp/munchies/internal/db"
	"mp/munchies/pkg/food"
	"mp/munchies/pkg/usda"
	"os"
	"strings"

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

type unitSet map[string]struct{}
type nutrientSet map[string]struct{}

// FIXME: extract, publish, unit test
func doIngest(usdaJsonPath string) error {
	foods, err := usda.MustRead(usdaJsonPath)
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error reading USDA JSON files at %q: %s.", usdaJsonPath, err), 1)
	}
	_ = os.RemoveAll(getDbPath())
	dB, err := db.New(getDbPath())
	if err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error creating database %q: %s.", getDbPath(), err), 1)
	}
	if err := dB.Migrate(); err != nil {
		log.Error().Err(err).Send()
		return cli.Exit(fmt.Sprintf("Error migrating database %q: %s.", getDbPath(), err), 1)
	}

	defer dB.Close()

	units := make(unitSet)
	nutrients := make(nutrientSet)

	for _, food := range foods {
		for _, fnut := range food.FoodNutrients {
			if err := doUnit(fnut, units, dB); err != nil {
				return err
			}
			if err := doNutrient(fnut, nutrients, dB); err != nil {
				return err
			}
		}
	}

	return nil
}

// FIXME: extract, publish, unit test
func doUnit(fnut food.FoodNutrient, units unitSet, dB *db.Database) error {
	name := fnut.Nutrient.UnitName
	if _, ok := units[name]; !ok {
		if err := (db.Unit{
			Name: strings.TrimSpace(name),
		}.WriteTo(dB)); err != nil {
			log.Error().Err(err).Send()
			return cli.Exit(fmt.Sprintf("Error writing unit %q to database %q: %s.", name, getDbPath(), err), 1)
		}
		units[name] = struct{}{}
	}

	return nil
}

// FIXME: extract, publish, unit test
func doNutrient(fnut food.FoodNutrient, nutrients nutrientSet, dB *db.Database) error {
	name := fnut.Nutrient.Name
	if _, ok := nutrients[name]; !ok {
		if err := (db.Nutrient{
			Name: strings.TrimSpace(name),
		}.WriteTo(dB)); err != nil {
			log.Error().Err(err).Send()
			return cli.Exit(fmt.Sprintf("Error writing unit %q to database %q: %s.", name, getDbPath(), err), 1)
		}
		nutrients[name] = struct{}{}
	}

	return nil
}
