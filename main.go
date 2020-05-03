package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:      "rssarchiver",
		Usage:     "Back up for RSS/ATOM feeds",
		ArgsUsage: "opml_file.xml",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return cli.ShowAppHelp(c)
			}
			return NewArchiver().UpdateFromOPML(c.Args().First())
		},
		Commands: []*cli.Command{
			{
				Name:  "summary",
				Usage: "Generate feed summary for a date",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Value:   time.Now().Format("2006-01-02"),
						Usage:   "Date for the summary",
					},
				},
				Action: func(c *cli.Context) error {
					return NewArchiver().GenerateSummary(c.String("date"))
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
