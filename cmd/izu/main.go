package main

import (
	"fmt"
	"os"

	"github.com/meir/izu/internal/izu"
	"github.com/urfave/cli/v2"
)

func main() {
	(&cli.App{
		Name:  "izu",
		Usage: "A unified hotkey config based on sxhkd.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to the configuration file",
			},
			&cli.StringFlag{
				Name:    "formatter",
				Aliases: []string{"f"},
				Usage:   "Path to the formatter lua file",
			},
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "Print verbose output",
			},
			&cli.StringFlag{
				Name:    "string",
				Aliases: []string{"s"},
				Usage:   "String to parse",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("formatter") == "" {
				return cli.Exit("formatter is required", 1)
			}

			formatter, err := izu.NewLuaFormatter(c.String("formatter"))
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if c.String("config") != "" {
				content, err := formatter.ParseFile(c.String("config"))
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
				fmt.Println(string(content))
			}

			if c.String("string") != "" {
				content, err := formatter.ParseString([]byte(c.String("string")))
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
				fmt.Println(string(content))
			}

			return nil
		},
	}).Run(os.Args)
}
