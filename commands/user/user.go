package user

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

type UserCommand struct {
	Settings *settings.BitAdminSettings
}

func (rc *UserCommand) GetCommand() cli.Command {

	userGrantCommand := &UserGrantCommand{
		Settings: rc.Settings,
		flags:    &UserGrantCommandFlags{},
	}

	return cli.Command{
		Name:  "user",
		Usage: "User opertations",
		Subcommands: []cli.Command{
			userGrantCommand.GetCommand(),
		},
	}
}
