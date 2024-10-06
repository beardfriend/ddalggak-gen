package main

import (
	"embed"
	"log"
	"os"

	"github.com/beardfriend/ddalggak-gen/internal"

	"github.com/urfave/cli/v2"
)

//go:embed all:template
var templates embed.FS

func main() {
	var workingDir, schemaName, modulePath string

	wid, err := os.Getwd()
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "workingDir",
				Value:       wid,
				Aliases:     []string{"w"},
				Usage:       "working directory",
				Destination: &workingDir,
			},
			&cli.StringFlag{
				Name:        "schemaName",
				Required:    true,
				Aliases:     []string{"s"},
				Usage:       "schema name",
				Destination: &schemaName,
			},
			&cli.StringFlag{
				Name:        "modulePath",
				Required:    true,
				Aliases:     []string{"m"},
				Usage:       "schema name",
				Destination: &modulePath,
			},
		},

		Name:  "ddalggak-gen",
		Usage: "boilerplate code generate tool",
		Action: func(*cli.Context) error {
			if err != nil {
				return err
			}
			field, err := internal.ParseEntity(workingDir, schemaName)
			if err != nil {
				return err
			}

			moduleName, err := internal.ParseProjectModuleName(workingDir)
			if err != nil {
				return err
			}
			repoTemp, err := templates.ReadFile("template/repo.tmpl")
			if err != nil {
				return err
			}
			err = internal.GenRepoFile(repoTemp, field, workingDir, modulePath, moduleName, schemaName)

			if err != nil {
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
