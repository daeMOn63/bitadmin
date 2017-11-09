// Package hooks hold actions on the Bitbucket hooks
package hooks

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// Command define base struct for repository subcommands and actions
type Command struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *Command) GetCommand() cli.Command {

	yaccHookCommand := YaccHookCommand{
		Settings: command.Settings,
	}

	return cli.Command{
		Name:  "hooks",
		Usage: "Hooks operations",
		Subcommands: []cli.Command{
			yaccHookCommand.GetCommand(),
		},
	}
}
