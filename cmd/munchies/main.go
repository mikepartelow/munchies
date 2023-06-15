package main

import (
	"fmt"
	"mp/munchies/internal/db"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	DB_PATH = "./munchies.sqlite"
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
			makeNutrientsCommand(),
			makeMealCommand(),
			makeUnitsCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}

func getDbPath() string {
	return DB_PATH
}

func getDb() (*db.Database, error) {
	return db.New(getDbPath())
}
