// Package user hold the actions on the Bitbucket users
package user

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// Command define base command for user actions
type Command struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (rc *Command) GetCommand() cli.Command {

	userGrantCommand := &GrantCommand{
		Settings: rc.Settings,
		flags:    &GrantCommandFlags{},
	}

	unsetPermissionsCommand := &UnsetPermissionsCommand{
		Settings: rc.Settings,
		flags:    &UnsetPermissionsCommandFlags{},
	}

	return cli.Command{
		Name:  "user",
		Usage: "User opertations",
		Subcommands: []cli.Command{
			userGrantCommand.GetCommand(),
			unsetPermissionsCommand.GetCommand(),
		},
	}
}
