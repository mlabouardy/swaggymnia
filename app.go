package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "swaggymnia"
	app.Usage = "Convert Insomnia to Swagger"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Mohamed Labouardy",
			Email: "mohamed@labouardy.com",
		},
	}
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "Load configuration from `FILE`",
				},
				cli.StringFlag{
					Name:  "insomnia, i",
					Usage: "Insomnia JSON `FILE`",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "yaml",
					Usage: "Output json|yaml",
				},
			},
			Usage: "Generate Swagger documentation",
			Action: func(c *cli.Context) error {
				if c.String("config") == "" || c.String("insomnia") == "" {
					return cli.NewExitError("config & insomnia flags are required", 1)
				}

				swagger := &Swagger{}
				swagger.Generate(c.String("insomnia"), c.String("config"), c.String("output"))

				return nil
			},
		},
	}

	app.Run(os.Args)
}
