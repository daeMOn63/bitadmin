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
func (command *Command) GetCommand() cli.Command {

	repositoryCreateCommand := &CreateCommand{
		Settings: command.Settings,
		flags:    &CreateCommandFlags{},
	}

	sonarCommand := &SonarCommand{
		Settings: command.Settings,
		flags:    &SonarCommandFlags{},
	}

	showPermissionsCommand := &ShowPermissionsCommand{
		Settings: command.Settings,
		flags:    &ShowPermissionsFlags{},
	}

	cloneSettingsCommand := &CloneSettingsCommand{
		Settings: command.Settings,
		flags:    &CloneSettingsCommandFlags{},
	}

	setBranchRestrictionCommand := &SetBranchRestrictionCommand{
		Settings: command.Settings,
		flags:    &SetBranchRestrictionCommandFlags{},
	}

	pullRequestSettingsCommand := &PullRequestSettingsCommand{
		Settings: command.Settings,
		flags:    &PullRequestSettingsCommandFlags{},
	}

	branchingModelCommand := &BranchingModelCommand{
		Settings: command.Settings,
		flags:    &BranchingModelCommandFlags{},
	}

	setDefaultReviewersCommand := &SetDefaultReviewersCommand{
		Settings: command.Settings,
		flags:    &SetDefaultReviewersCommandFlags{},
	}

	moveCommand := &MoveCommand{
		Settings: command.Settings,
		flags:    &MoveCommandFlags{},
	}

	return cli.Command{
		Name:  "repository",
		Usage: "Repository operations",
		Subcommands: []cli.Command{
			repositoryCreateCommand.GetCommand(),
			sonarCommand.GetCommand(),
			showPermissionsCommand.GetCommand(),
			cloneSettingsCommand.GetCommand(),
			setBranchRestrictionCommand.GetCommand(),
			pullRequestSettingsCommand.GetCommand(),
			branchingModelCommand.GetCommand(),
			setDefaultReviewersCommand.GetCommand(),
			moveCommand.GetCommand(),
		},
	}
}
