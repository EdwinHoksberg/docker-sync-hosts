package utility

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	MARKER_BEGIN = "#docker-sync-hosts begin\n"
	MARKER_END   = "#docker-sync-hosts end\n"
)

type hostsSync struct {
	Extension string
	HostsFile *hostsFile
}

func NewHostsSync(extension string) *hostsSync {
	return &hostsSync{
		Extension: extension,
		HostsFile: NewHostsFile(),
	}
}

func (h *hostsSync) Write(containerMap map[string][]string) error {
	// Create a regex to find our markers in the hosts file.
	regex, err := regexp.Compile(fmt.Sprintf("(?s)%s(.*)%s", MARKER_BEGIN, MARKER_END))
	if err != nil {
		return err
	}

	// Read the contents of the hosts file.
	hostsContent, err := h.HostsFile.Read()
	if err != nil {
		return err
	}

	// Empty our markers and everything in between to prepare for writing the new hosts.
	replacedHosts := regex.ReplaceAllString(hostsContent, "")

	// Create a buffer for writing the new host lines.
	result := bytes.NewBuffer([]byte(replacedHosts))
	result.WriteString(MARKER_BEGIN)

	// Loop through each container, and add the ip and aliases to the buffer.
	for ip, aliases := range containerMap {
		var hostnames bytes.Buffer
		for i, name := range aliases {
			hostnames.WriteString(name)

			// Check if the alias already has a extension which the user provided.
			if !strings.HasSuffix(name, h.Extension) {
				hostnames.WriteString(h.Extension)
			}

			// Only write a splitting space character if it is not the last alias provided.
			if len(aliases)-1 != i {
				hostnames.WriteRune(' ')
			}
		}

		result.WriteString(fmt.Sprintf("%s %s\n", ip, hostnames.String()))
	}

	// Finally, write a marker to signal the end of our block of host lines.
	result.WriteString(MARKER_END)

	// After we looped over every container, write the contents back to the hosts file.
	if err := h.HostsFile.Write(result.String(), false); err != nil {
		return err
	}

	return nil
}

func (h *hostsSync) AddEntry(ip string, aliases []string) error {
	// Read the contents of the hosts file.
	hostsContent, err := h.HostsFile.Read()
	if err != nil {
		return err
	}

	// Parse the current host lines from the hosts file.
	containerMap, err := h.parseBlock(hostsContent)
	if err != nil {
		return err
	}

	// Add the new container to the host lines.
	containerMap[ip] = aliases

	// Write the new list to the hosts file.
	if err := h.Write(containerMap); err != nil {
		return err
	}

	return nil
}

func (h *hostsSync) RemoveEntry(ip string) error {
	// Read the contents of the hosts file.
	hostsContent, err := h.HostsFile.Read()
	if err != nil {
		return err
	}

	// Parse the current host lines from the hosts file.
	containerMap, err := h.parseBlock(hostsContent)
	if err != nil {
		return err
	}

	// Check if the ip exists in the host lines, and remove it if it does.
	_, exists := containerMap[ip]
	if exists {
		delete(containerMap, ip)
	}

	// Write the new list to the hosts file.
	if err := h.Write(containerMap); err != nil {
		return err
	}

	return nil
}

func (h *hostsSync) parseBlock(hostsContent string) (map[string][]string, error) {
	// Create a regex to find our markers in the hosts file.
	regex, err := regexp.Compile(fmt.Sprintf("(?s)%s(.*)%s", MARKER_BEGIN, MARKER_END))
	if err != nil {
		return nil, err
	}

	// Initialize a variable for holding the host lines.
	containerMap := make(map[string][]string)

	// Execute the regex to find our host lines block.
	matches := regex.FindStringSubmatch(hostsContent)

	// If no match and group was found, return an empty map.
	if len(matches) < 2 {
		return containerMap, nil
	}

	// Split the host lines in multiple chuncks.
	hostLines := strings.Split(matches[1], "\n")

	// Loop through each host line.
	for _, hostLine := range hostLines {
		// Trim any newlines and spaces from the hostline
		hostLine = strings.Trim(hostLine, "\r\n ")

		// If it was empty after trimming, we can skip it.
		if len(hostLine) == 0 {
			continue
		}

		// The first part of a host line is always the ip, all the next values are the container aliases.
		splittedMatch := strings.SplitN(hostLine, " ", 2)

		ip := splittedMatch[0]
		aliases := splittedMatch[1:]

		// Add the container to our map.
		containerMap[ip] = aliases
	}

	return containerMap, nil
}
