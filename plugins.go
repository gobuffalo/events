package events

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/pkg/errors"
)

// LoadPlugins will add listeners for any plugins that support "events"
func LoadPlugins() error {
	plugs, err := plugins.Available()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, cmds := range plugs {
		for _, c := range cmds {
			if c.BuffaloCommand != "events" {
				continue
			}
			NamedListen(fmt.Sprintf("[PLUGIN] %s %s", c.Binary, c.Name), func(e Event) {
				b, err := json.Marshal(e)
				if err != nil {
					fmt.Println("error trying to marshal event", e, err)
					return
				}
				cmd := exec.Command(c.Binary, c.UseCommand, string(b))
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				if err := cmd.Run(); err != nil {
					fmt.Println("error trying to send event", strings.Join(cmd.Args, " "), err)
				}
			})
		}

	}
	return nil
}

func init() {
	LoadPlugins()
}
