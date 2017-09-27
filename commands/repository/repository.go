package repository

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

type RepositoryCommand struct {
	Settings *settings.BitAdminSettings
}

func (rc *RepositoryCommand) GetCommand() cli.Command {

	fileCache := rc.Settings.GetFileCache()

	repositoryCreateCommand := &RepositoryCreateCommand{
		Settings: rc.Settings,
		flags:    &RepositoryCreateCommandFlags{},
	}

	enableSonarCleanupCommand := &EnableSonarCleanupCommand{
		Settings: rc.Settings,
	}

	return cli.Command{
		Name:  "repository",
		Usage: "Repository opertations",
		Subcommands: []cli.Command{
			repositoryCreateCommand.GetCommand(fileCache),
			enableSonarCleanupCommand.GetCommand(),
		},
	}
}
