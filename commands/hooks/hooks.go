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

	eolHookCommand := EolHookCommand{
		Settings: command.Settings,
	}

	pubHookCommand := PubHookCommand{
		Settings: command.Settings,
	}

	yaccHookCommand := YaccHookCommand{
		Settings: command.Settings,
	}

	rfpHookCommand := RfpHookCommand{
		Settings: command.Settings,
	}

	return cli.Command{
		Name:  "hooks",
		Usage: "Hooks operations",
		Subcommands: []cli.Command{
			eolHookCommand.GetCommand(),
			pubHookCommand.GetCommand(),
			yaccHookCommand.GetCommand(),
			rfpHookCommand.GetCommand(),
		},
	}
}
