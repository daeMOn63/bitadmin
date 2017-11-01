package group

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

type GroupCommand struct {
	Settings *settings.BitAdminSettings
}

func (rc *GroupCommand) GetCommand() cli.Command {

	groupGrantCommand := &GroupGrantCommand{
		Settings: rc.Settings,
		flags:    &GroupGrantCommandFlags{},
	}

	return cli.Command{
		Name:  "group",
		Usage: "Group opertations",
		Subcommands: []cli.Command{
			groupGrantCommand.GetCommand(),
		},
	}
}
