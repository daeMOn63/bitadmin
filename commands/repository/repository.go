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
	fileCache.Load()

	repositoryCreateCommand := &RepositoryCreateCommand{
		Settings: rc.Settings,
		flags:    &RepositoryCreateCommandFlags{},
	}

	return cli.Command{
		Name:  "repository",
		Usage: "Repository opertations",
		Subcommands: []cli.Command{
			repositoryCreateCommand.GetCommand(fileCache),
		},
	}
}
