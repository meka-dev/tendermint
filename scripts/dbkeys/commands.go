package main

import (
	"fmt"
	"os"

	"github.com/creachadair/command"
	"github.com/tendermint/tendermint/scripts/dbkeys/config"

	tmdb "github.com/tendermint/tm-db"
)

var cmdList = &command.C{
	Name: "list",
	Help: "List all the keys in the database.",

	Run: func(env *command.Env, args []string) error {
		cfg := env.Config.(*config.Settings)
		return cfg.WithDB(func(db tmdb.DB) error {
			it, err := db.Iterator(nil, nil)
			if err != nil {
				return fmt.Errorf("create iterator: %w", err)
			}

			for it.Valid() {
				fmt.Printf("%#q\n", string(it.Key()))
				it.Next()
			}
			return it.Close()
		})
	},
}

var cmdStats = &command.C{
	Name: "stats",
	Help: "Print summary stats for the database.",

	Run: func(env *command.Env, args []string) error {
		cfg := env.Config.(*config.Settings)
		return cfg.WithDB(func(db tmdb.DB) error {
			// The DB interface has a Stats method, but it isn't always implemented.
			it, err := db.Iterator(nil, nil)
			if err != nil {
				return fmt.Errorf("create iterator: %w", err)
			}

			stats := new(config.Histogram)
			for it.Valid() {
				stats.AddSample(len(it.Value()))
				it.Next()
			}

			stats.WriteTo(os.Stdout)
			return it.Close()
		})
	},
}
