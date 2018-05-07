package main

import (
	"fmt"
	"github.com/edwinhoksberg/docker-sync-hosts/command"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	Name    string
	Version string
)

func main() {
	command := new(command.Command)

	app := cli.NewApp()

	app.Name = Name
	app.Usage = ""
	app.HelpName = Name
	app.Version = Version
	app.Authors = []cli.Author{
		{
			Name:  "Edwin Hoksberg",
			Email: "mail@edwinhoksberg.nl",
		},
	}
	app.UsageText = fmt.Sprintf(`A simple cli application to keep your hosts file up-to-date with running docker containers.

   Homepage: https://github.com/edwinhoksberg/docker-sync-hosts`)

	app.Commands = []cli.Command{
		{
			Name:   "sync",
			Usage:  "Sync the hosts file with currently running docker containers",
			Action: command.Sync,
		},
		{
			Name:   "daemon",
			Usage:  "Start a daemon to sync your hosts file when a container is started or stopped",
			Action: command.Daemon,
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket",
			Usage: "socket endpoint for docker",
			Value: "unix:///var/run/docker.sock",
		},
		cli.StringFlag{
			Name:  "extension",
			Usage: "the hostname extension to use",
			Value: ".test",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable debug logging",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "completely disable logging",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error()
	}
}
