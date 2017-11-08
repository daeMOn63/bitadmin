// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// PullRequestSettingsCommand define base struct for SetBranchRestriction actions
type PullRequestSettingsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *PullRequestSettingsCommandFlags
}

// PullRequestSettingsCommandFlags hold flag values for the PullRequestSettingsCommand
type PullRequestSettingsCommandFlags struct {
	project                  string
	repository               string
	requiredAllApprovers     bool
	requiredAllTaskComplete  bool
	unapproveOnUpdate        bool
	requiredApprovers        uint
	requiredSuccessfulBuilds uint
}

// GetCommand provide a ready to use cli.Command
func (command *PullRequestSettingsCommand) GetCommand(fileCache *helper.FileCache) cli.Command {
	return cli.Command{
		Name:   "set-pr-settings",
		Usage:  "Set pull request settings on given repository",
		Action: command.SetPullRequestSettingsAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<project_key>` where the repository will be created",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository_name>` to create",
				Destination: &command.flags.repository,
			},
			cli.BoolFlag{
				Name:        "requiredAllApprovers",
				Usage:       "Tell when all reviewers must have approved to allow merge",
				Destination: &command.flags.requiredAllApprovers,
			},
			cli.BoolFlag{
				Name:        "requiredAllTaskComplete",
				Usage:       "Tell when all tasks must have completed to allow merge",
				Destination: &command.flags.requiredAllTaskComplete,
			},
			cli.BoolFlag{
				Name:        "unapproveOnUpdate",
				Usage:       "Tell if all approvals should be removed when the pull request get updated",
				Destination: &command.flags.unapproveOnUpdate,
			},
			cli.UintFlag{
				Name:        "requiredApprovers",
				Usage:       "`<requiredApprovers>` set the minimum number of approval required to allow merge",
				Destination: &command.flags.requiredApprovers,
			},
			cli.UintFlag{
				Name:        "requiredSuccessfulBuilds",
				Usage:       "`<requiredSuccessfulBuilds>` set the minimum number of successful builds required to allow merge",
				Destination: &command.flags.requiredSuccessfulBuilds,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, fileCache)
		},
	}
}

// SetPullRequestSettingsAction allow to set the pull request settings on given repository.
// This will keep the existing MergeConfig and doesn't allow to update this for now.
func (command *PullRequestSettingsCommand) SetPullRequestSettingsAction(context *cli.Context) error {
	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	pullRequestSettings, err := client.GetPullRequestSettings(command.flags.project, command.flags.repository)

	pullRequestSettings.RequiredAllApprovers = command.flags.requiredAllApprovers
	pullRequestSettings.RequiredAllTasksComplete = command.flags.requiredAllTaskComplete
	pullRequestSettings.RequiredApprovers = command.flags.requiredApprovers
	pullRequestSettings.RequiredSuccessfulBuilds = command.flags.requiredSuccessfulBuilds
	pullRequestSettings.UnapproveOnUpdate = command.flags.unapproveOnUpdate

	err = client.SetPullRequestSettings(command.flags.project, command.flags.repository, pullRequestSettings)
	if err != nil {
		return err
	}

	fmt.Printf("Pull request settings successfully set on %s/%s\n", command.flags.project, command.flags.repository)

	return nil

}
