package command

import (
	"github.com/edwinhoksberg/docker-sync-hosts/utility"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func (c *Command) Sync(cli *cli.Context) {
	log := c.GetLogger(cli)

	// Build and connect to the docker remote api.
	client, err := c.GetDockerSocketClient(cli)
	if err != nil {
		log.WithError(err).Error("Could not connect to docker socket")
		return
	}

	// List all available containers.
	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.WithError(err).Error("Failed to retrieve list of containers")
		return
	}

	// Check if any containers where running.
	if len(containers) == 0 {
		log.Info("No containers found, quiting...")
		return
	}

	// Variable for storing all found containers and their ip addresses
	containerMap := make(map[string][]string)

	// Loop through each container and
	for _, container := range containers {
		containerObj, err := client.InspectContainer(container.ID)
		if err != nil {
			log.WithError(err).WithField("container", container.ID).Error("Failed to inspect container")
			return
		}

		// Loop through all networks, and add the found ip/aliases to the map variable.
		for _, network := range containerObj.NetworkSettings.Networks {
			// Don't add container ip without aliases to map.
			if len(network.Aliases) == 0 {
				continue
			}

			containerMap[network.IPAddress] = network.Aliases

			log.WithFields(logrus.Fields{
				"aliases":   network.Aliases,
				"container": containerObj.Name,
				"ip":        network.IPAddress,
			}).Debug("Found container network")
		}
	}

	// Initialize the hosts syncing class.
	hostsSync := utility.NewHostsSync(cli.GlobalString("extension"))

	// Write the newly generated hosts file.
	if err := hostsSync.Write(containerMap); err != nil {
		log.WithError(err).Error("Failed to write to hosts file")
		return
	}
}
