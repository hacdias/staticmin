package runner

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/users"
)

// Runner is a commands runner.
type Runner struct {
	Enabled bool
	*settings.Settings
}

// RunBeforeHook and RunAfterHook trigger the command runner at their respective times.
func (r *Runner) RunBeforeHook(evt, path, dst string, user *users.User) error {
	path = user.FullPath(path)
	dst = user.FullPath(dst)
	if r.Enabled {
		if val, ok := r.Commands["before_"+evt]; ok {
			for _, command := range val {
				err := r.exec(command, "before_"+evt, path, dst, user)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *Runner) RunAfterHook(evt, path, dst string, user *users.User) error {
	path = user.FullPath(path)
	dst = user.FullPath(dst)
	if r.Enabled {
		if val, ok := r.Commands["after_"+evt]; ok {
			for _, command := range val {
				err := r.exec(command, "after_"+evt, path, dst, user)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// RunHook runs the hooks for the before and after event.
func (r *Runner) RunHook(fn func() error, evt, path, dst string, user *users.User) error {
	path = user.FullPath(path)
	dst = user.FullPath(dst)

	if r.Enabled {
		if val, ok := r.Commands["before_"+evt]; ok {
			for _, command := range val {
				err := r.exec(command, "before_"+evt, path, dst, user)
				if err != nil {
					return err
				}
			}
		}
	}

	err := fn()
	if err != nil {
		return err
	}

	if r.Enabled {
		if val, ok := r.Commands["after_"+evt]; ok {
			for _, command := range val {
				err := r.exec(command, "after_"+evt, path, dst, user)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *Runner) exec(raw, evt, path, dst string, user *users.User) error {
	blocking := true

	if strings.HasSuffix(raw, "&") {
		blocking = false
		raw = strings.TrimSpace(strings.TrimSuffix(raw, "&"))
	}

	command, err := ParseCommand(r.Settings, raw)
	if err != nil {
		return err
	}

	envMapping := func(key string) string {
		switch key {
		case "FILE":
			return path
		case "SCOPE":
			return user.Scope
		case "TRIGGER":
			return evt
		case "USERNAME":
			return user.Username
		case "DESTINATION":
			return dst
		default:
			return os.Getenv(key)
		}
	}
	for i, arg := range command {
		if i == 0 {
			continue
		}

		command[i] = os.Expand(arg, envMapping)
	}

	cmd := exec.Command(command[0], command[1:]...) //nolint:gosec
	cmd.Env = append(os.Environ(), fmt.Sprintf("FILE=%s", path))
	cmd.Env = append(cmd.Env, fmt.Sprintf("SCOPE=%s", user.Scope)) //nolint:gocritic
	cmd.Env = append(cmd.Env, fmt.Sprintf("TRIGGER=%s", evt))
	cmd.Env = append(cmd.Env, fmt.Sprintf("USERNAME=%s", user.Username))
	cmd.Env = append(cmd.Env, fmt.Sprintf("DESTINATION=%s", dst))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if !blocking {
		log.Printf("[INFO] Nonblocking Command: \"%s\"", strings.Join(command, " "))
		defer func() {
			go func() {
				err := cmd.Wait()
				if err != nil {
					log.Printf("[INFO] Nonblocking Command \"%s\" failed: %s", strings.Join(command, " "), err)
				}
			}()
		}()
		return cmd.Start()
	}

	log.Printf("[INFO] Blocking Command: \"%s\"", strings.Join(command, " "))
	return cmd.Run()
}
