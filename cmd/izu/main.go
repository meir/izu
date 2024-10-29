package main

import (
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/meir/izu/internal/luaformatter"
	"github.com/meir/izu/internal/parser"
	"github.com/meir/izu/pkg/izu"
	"github.com/phsym/console-slog"
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
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"S"},
				Usage:   "Silent output, does not output any logs or errors unless when panicking",
			},
			&cli.StringFlag{
				Name:    "string",
				Aliases: []string{"s"},
				Usage:   "String to parse",
			},
		},
		Action: func(c *cli.Context) error {
			level := slog.LevelInfo
			if c.Bool("verbose") {
				level = slog.LevelDebug
			} else if c.Bool("silent") {
				level = math.MaxInt32 // Never log, except for the output using fmt
			}
			slog.SetDefault(slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
				Level: level,
			})))

			if c.Bool("version") {
				slog.Info("Izu Version " + izu.GetVersion())
				return nil
			}

			input := []byte(c.String("string"))
			if c.String("config") != "" {
				content, err := os.ReadFile(c.String("config"))
				if err != nil {
					slog.Error("Failed to read config file: " + err.Error())
					return cli.Exit("", 1)
				}
				input = content
			}

			hotkeys, err := parser.Parse([]byte(input))
			if err != nil {
				slog.Error("Failed to parse hotkeys: " + err.Error())
				return cli.Exit("", 1)
			}

			formatter, err := luaformatter.NewFormatter(c.String("formatter"))
			if err != nil {
				slog.Error("Failed to create formatter: " + err.Error())
				return cli.Exit("", 1)
			}

			lines, err := formatter.Format(hotkeys)
			if err != nil {
				slog.Error("Failed to format hotkeys: " + err.Error())
				return cli.Exit("", 1)
			}

			for _, line := range lines {
				fmt.Println(line)
			}

			return nil
		},
	}).Run(os.Args)
}
