package command

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/urfave/cli"
)

type Command struct{}

func (c Command) GetLogger(cli *cli.Context) *logrus.Logger {
	if cli.GlobalBool("quiet") {
		logger, _ := test.NewNullLogger()
		return logger
	}

	// Set the default output formatter
	format := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	}

	logrus.SetFormatter(format)

	// If the verbose flag was enabled, enable debug logging
	if cli.GlobalBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return logrus.StandardLogger()
}

func (c Command) GetDockerSocketClient(cli *cli.Context) (*docker.Client, error) {
	logrus.Debugf("Connecting to docker: %s", cli.GlobalString("socket"))

	// Connect to the beanstalkd server.
	client, err := docker.NewClient(cli.GlobalString("socket"))
	if err != nil {
		return nil, err
	}

	logrus.Debug("Succesfully connected to docker socket")

	return client, nil
}
