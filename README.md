# docker-sync-hosts [![Build Status](https://travis-ci.org/EdwinHoksberg/docker-sync-hosts.svg?branch=master)](https://travis-ci.org/EdwinHoksberg/docker-sync-hosts)
A simple cli application to keep your hosts file up-to-date with running docker containers.

## Download
You can find the latest release including binaries [here](https://github.com/EdwinHoksberg/docker-sync-hosts/releases/latest).

## Usage
```
NAME:
   docker-sync-hosts

USAGE:
   A simple cli application to keep your hosts file up-to-date with running docker containers.

   Homepage: https://github.com/edwinhoksberg/docker-sync-hosts

VERSION:
    []

AUTHOR:
   Edwin Hoksberg <mail@edwinhoksberg.nl>

COMMANDS:
     sync     Sync the hosts file with currently running docker containers
     daemon   Start a daemon to sync your hosts file when a container is started or stopped
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --socket value     socket endpoint for docker (default: "unix:///var/run/docker.sock")
   --extension value  the hostname extension to use (default: ".test")
   --verbose          enable debug logging
   --quiet            completely disable logging
   --help, -h         show help
   --version, -v      print the version
```

## Development
This program is written in [Go](https://golang.org/), using these dependencies:
- [sirupsen/logrus](https://github.com/sirupsen/logrus) - Logrus is a structured logger for Go (golang), completely API compatible with the standard library logger.
- [urfave/cli](https://github.com/urfave/cli) - cli is a simple, fast, and fun package for building command line apps in Go.
- [fsouza/go-dockerclient](https://github.com/fsouza/go-dockerclient) - Go client for the Docker remote API.

## License
[MIT](LICENSE.md)
