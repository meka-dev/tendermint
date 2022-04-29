package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/creachadair/command"
	"github.com/tendermint/tendermint/scripts/dbkeys/config"
)

var (
	dbSpec   = flag.String("db", "", `Database spec ("backend:path", required)`)
	statBase = flag.Int("base", 4, "Histogram sample base")
	statBars = flag.Int("bar", 20, "Histogram bar length")
)

func main() {
	var cfg config.Settings
	root := &command.C{
		Name:  filepath.Base(os.Args[0]),
		Usage: "-db spec command [args...]",
		Help:  "Manipulate the contents of a Tendermint database file.",

		SetFlags: func(_ *command.Env, fs *flag.FlagSet) {
			fs.StringVar(&cfg.Spec, "db", "", `Database spec (backend:path, required)`)
		},

		Commands: []*command.C{
			cmdList,
			cmdStats,
			command.HelpCommand(nil),
		},
	}
	command.RunOrFail(root.NewEnv(&cfg), os.Args[1:])
}
