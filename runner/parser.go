package runner

import (
	"os/exec"

	"github.com/filebrowser/filebrowser/v2/settings"
)

// ParseCommand parses the command taking in account if the current
// instance uses a shell to run the commands or just calls the binary
// directly.
func ParseCommand(s *settings.Settings, raw string) ([]string, error) {
	var command []string

	if len(s.Shell) == 0 || s.Shell[0] == "" {
		cmd, args, err := SplitCommandAndArgs(raw)
		if err != nil {
			return nil, err
		}

		_, err = exec.LookPath(cmd)
		if err != nil {
			return nil, err
		}

		command = append(command, cmd)
		command = append(command, args...)
	} else {
		cmd, args, err := SplitCommandAndArgs(s.Shell[0])
		if err != nil {
			return nil, err
		}

		_, err = exec.LookPath(cmd)
		if err != nil {
			return nil, err
		}

		command = append(command, cmd)
		command = append(command, args...)
		command = append(command, raw)
	}

	return command, nil
}
