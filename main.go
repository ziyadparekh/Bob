package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ziyadparekh/bob/builder"
)

func main() {
	app := cli.NewApp()
	app.Name = "bob"
	app.Usage = "Build Services to Deathstar"
	app.Author = "Ziyad Parekh"
	app.Commands = []cli.Command{
		{
			Name:  "build",
			Usage: "Build a service to deathstar!",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "branch, b",
					Value: "",
					Usage: "Branch to build",
				},
				cli.StringFlag{
					Name:  "service, s",
					Value: "",
					Usage: "Pick a service to build or say '-s all' to build all",
				},
				cli.StringFlag{
					Name:  "env, e",
					Value: "qa",
					Usage: "Pick between development or qa",
				},
				cli.BoolFlag{
					Name:  "open, o",
					Usage: "Open the job in your browser",
				},
			},
			Action: func(c *cli.Context) {
				b, _ := builder.NewBob()
				if c.String("service") == "all" {
					b.BuildAllServices(c.String("branch"), c.Bool("open"))
				} else if c.String("service") == "adminv2" {
					b.BuildClient(c.String("service"), c.String("branch"), c.String("env"), c.Bool("open"))
				} else {
					b.BuildService(c.String("service"), c.String("branch"), c.Bool("open"))
				}
			},
		},
		{
			Name:  "list",
			Usage: "List all the services that are buildable",
			Action: func(c *cli.Context) {
				fmt.Println("=========SERVICES=========")
				for k := range builder.JenkinsServices {
					str := fmt.Sprintf("-----> %s", k)
					fmt.Println(str)
				}
				fmt.Println("=========CLIENTS=========")
				for k := range builder.JenkinsClients {
					str := fmt.Sprintf("-----> %s", k)
					fmt.Println(str)
				}

			},
		},
	}

	app.Run(os.Args)
}
