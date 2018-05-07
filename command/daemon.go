package command

import (
	"github.com/edwinhoksberg/docker-sync-hosts/utility"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

type daemon struct {
	cli           *cli.Context
	log           *logrus.Logger
	client        *docker.Client
	eventListener chan *docker.APIEvents
}

func (c *Command) Daemon(cli *cli.Context) {
	log := c.GetLogger(cli)

	// Build and connect to the docker remote api.
	client, err := c.GetDockerSocketClient(cli)
	if err != nil {
		log.WithError(err).Error("Could not connect to docker socket")
		return
	}

	// Do a full sync of the containers before we begin.
	c.Sync(cli)

	// Create and attach event listener.
	eventListener := make(chan *docker.APIEvents)
	if err := client.AddEventListener(eventListener); err != nil {
		log.WithError(err).Error("Failed to attach docker event listener")
		return
	}

	// Create a simple holder for passing around these variables.
	daemon := &daemon{
		cli:           cli,
		log:           log,
		client:        client,
		eventListener: eventListener,
	}

	// Start a loop to catch any signal events, like SIGTERM.
	catchSignals(daemon)

	log.Info("Listening for container events...")

	// Infinite loop for retrieving events.
	for {
		select {
		case msg := <-eventListener:
			processEvent(daemon, msg)
		}
	}
}

func processEvent(daemon *daemon, event *docker.APIEvents) {
	// Check if event was a container start or stop(kill) event.
	if event.Action != "start" && event.Action != "kill" {
		return
	}

	// Lookup the container which triggered the event.
	container, err := daemon.client.InspectContainer(event.Actor.ID)
	if err != nil {
		daemon.log.WithError(err).WithField("container", event.Actor.ID).Error("Failed to inspect container")
		return
	}

	var ip string
	var aliases []string

	// Loop through all networks, and add the found ip/aliases to the map variable.
	for _, network := range container.NetworkSettings.Networks {
		ip = network.IPAddress
		aliases = network.Aliases

		daemon.log.WithFields(logrus.Fields{
			"aliases":   network.Aliases,
			"container": container.Name,
			"ip":        network.IPAddress,
		}).Debug("Found container network")
	}

	// Don't add container ip without any aliases.
	if len(aliases) == 0 {
		return
	}

	// Initialize the hosts syncing class.
	hostsSync := utility.NewHostsSync(daemon.cli.GlobalString("extension"))

	// When a container is starting up, add the ip and aliases to the hosts file.
	if event.Action == "start" {
		daemon.log.WithField("container", container.Name).Info("Adding container...")

		if err := hostsSync.AddEntry(ip, aliases); err != nil {
			daemon.log.WithError(err).Error("Failed to remove container from hosts file")
		}
	}

	// When a container is being stopped, remove the ip and aliases from the hosts file.
	if event.Action == "kill" {
		daemon.log.WithField("container", container.Name).Info("Removing container...")

		if err := hostsSync.RemoveEntry(ip); err != nil {
			daemon.log.WithError(err).Error("Failed to remove container from hosts file")
		}
	}
}

func catchSignals(daemon *daemon) {
	signal_channel := make(chan os.Signal, 1)

	// Wait for a Interrupt or sigterm signal, and handle it.
	signal.Notify(signal_channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			recv_signal := <-signal_channel

			// We only catch interrupt events, like CTRL^C
			switch recv_signal {
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				daemon.log.Info("Removing event listener and shutting down...")

				daemon.client.RemoveEventListener(daemon.eventListener)
				close(daemon.eventListener)

				os.Exit(0)
			}
		}
	}()
}
