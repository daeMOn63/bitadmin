// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// Command define base struct for repository subcommands and actions
type Command struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (rc *Command) GetCommand() cli.Command {

	fileCache := rc.Settings.GetFileCache()

	repositoryCreateCommand := &CreateCommand{
		Settings: rc.Settings,
		flags:    &CreateCommandFlags{},
	}

	enableSonarCleanupCommand := &EnableSonarCleanupCommand{
		Settings: rc.Settings,
	}

	showPermissionsCommand := &ShowPermissionsCommand{
		Settings: rc.Settings,
		flags:    &ShowPermissionsFlags{},
	}

	cloneSettingsCommand := &CloneSettingsCommand{
		Settings: rc.Settings,
		flags:    &CloneSettingsCommandFlags{},
	}

	return cli.Command{
		Name:  "repository",
		Usage: "Repository opertations",
		Subcommands: []cli.Command{
			repositoryCreateCommand.GetCommand(fileCache),
			enableSonarCleanupCommand.GetCommand(),
			showPermissionsCommand.GetCommand(fileCache),
			cloneSettingsCommand.GetCommand(fileCache),
		},
	}
}
