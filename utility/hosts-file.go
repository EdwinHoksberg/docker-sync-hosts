package utility

import (
	"bytes"
	"io"
	"os"
	"runtime"
)

const (
	LINUX_HOSTS_FILE   = "/etc/hosts"
	WINDOWS_HOSTS_FILE = "C:\\Windows\\System32\\drivers\\etc\\hosts"
)

type hostsFile struct{}

func NewHostsFile() *hostsFile {
	return &hostsFile{}
}

func (h *hostsFile) Read() (string, error) {
	// Create a new buffer for reading the file contents into.
	result := bytes.NewBuffer(nil)

	// Open the hosts file in the correct file mode.
	f, err := os.OpenFile(h.getHostsFile(), os.O_RDONLY, h.getHostsFileMode())
	if err != nil {
		return "", err
	}

	// Copy the file contents buffer to our own buffer.
	io.Copy(result, f)

	// Close the file.
	f.Close()

	return result.String(), nil
}

func (h *hostsFile) Truncate(toSize int) error {
	// Open the hosts file in the correct file mode.
	f, err := os.OpenFile(h.getHostsFile(), os.O_WRONLY, h.getHostsFileMode())
	if err != nil {
		return err
	}

	// Truncate the file down to the specifed size.
	if err = f.Truncate(int64(toSize)); err != nil {
		return err
	}

	// Make sure the file marker is set at the new end of the file.
	f.Seek(int64(toSize), 0)

	// Close the file.
	f.Close()

	return nil
}

func (h *hostsFile) Write(hosts string, append bool) error {
	// Open the hosts file in the correct file mode.
	f, err := os.OpenFile(h.getHostsFile(), os.O_WRONLY, h.getHostsFileMode())
	if err != nil {
		return err
	}

	// If we are here to overwrite the file, empty the file first.
	if !append {
		if err = h.Truncate(0); err != nil {
			return err
		}
	}

	// Write the string to the file.
	_, err = f.WriteString(hosts)
	if err != nil {
		return err
	}

	// Close the file.
	f.Close()

	return nil
}

func (h *hostsFile) getHostsFile() string {
	// Determine the hosts file location.
	if runtime.GOOS == "windows" {
		return WINDOWS_HOSTS_FILE
	}

	return LINUX_HOSTS_FILE
}

func (h *hostsFile) getHostsFileMode() os.FileMode {
	// Find the file mode for the hosts file.
	stat, _ := os.Stat(h.getHostsFile())

	return stat.Mode()
}
