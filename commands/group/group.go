// Package group provides actions for interacting with Bitbucket groups.
package group

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// Command define the base struct for providing group actions
type Command struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (rc *Command) GetCommand() cli.Command {

	grantCommand := &GrantCommand{
		Settings: rc.Settings,
		flags:    &GrantCommandFlags{},
	}

	unsetPermissionsCommand := &UnsetPermissionsCommand{
		Settings: rc.Settings,
		flags:    &UnsetPermissionsCommandFlags{},
	}

	return cli.Command{
		Name:  "group",
		Usage: "Group opertations",
		Subcommands: []cli.Command{
			grantCommand.GetCommand(),
			unsetPermissionsCommand.GetCommand(),
		},
	}
}
