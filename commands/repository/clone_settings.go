// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// CloneSettingsCommand define base struct for the clone settings actions
type CloneSettingsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *CloneSettingsCommandFlags
}

// CloneSettingsCommandFlags hold flag values for the CloneSettingsCommand
type CloneSettingsCommandFlags struct {
	sourceRepository    string
	targetRepository    string
	userPermissions     bool
	groupPermissions    bool
	branchRestrictions  bool
	pullRequestSettings bool
}

// GetCommand provide a ready to use cli.Command
func (command *CloneSettingsCommand) GetCommand(fileCache *helper.FileCache) cli.Command {
	return cli.Command{
		Name:   "clone-settings",
		Usage:  "Clone various settings from a repository to another",
		Action: command.CloneSettingsAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "sourceRepository",
				Usage:       "The `<sourceRepository>` to read the settings from",
				Destination: &command.flags.sourceRepository,
			},
			cli.StringFlag{
				Name:        "targetRepository",
				Usage:       "The `<targetRepository>` to copy the settings to",
				Destination: &command.flags.targetRepository,
			},
			cli.BoolFlag{
				Name:        "userPermissions",
				Usage:       "Copy user permissions",
				Destination: &command.flags.userPermissions,
			},
			cli.BoolFlag{
				Name:        "groupPermissions",
				Usage:       "Copy group permissions",
				Destination: &command.flags.groupPermissions,
			},
			cli.BoolFlag{
				Name:        "branchRestrictions",
				Usage:       "Copy branch restrictions",
				Destination: &command.flags.branchRestrictions,
			},
			cli.BoolFlag{
				Name:        "pullRequestSettings",
				Usage:       "Copy pull-request settings",
				Destination: &command.flags.pullRequestSettings,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, fileCache)
		},
	}
}

// CloneSettingsAction provide logic allowing to copy repository settings from one to another.
// Thoses settings include user / group permissions, and branch restrictions.
func (command *CloneSettingsCommand) CloneSettingsAction(context *cli.Context) error {
	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	fileCache := command.Settings.GetFileCache()

	sourceRepo, err := fileCache.SearchRepositorySlug(command.flags.sourceRepository)
	if err != nil {
		return err
	}

	targetRepo, err := fileCache.SearchRepositorySlug(command.flags.targetRepository)
	if err != nil {
		return err
	}

	if command.flags.userPermissions == true {
		err := client.CloneRepositoryUserPermissions(sourceRepo.Project.Key, sourceRepo.Slug, targetRepo.Project.Key, targetRepo.Slug)
		if err != nil {
			return err
		}
		fmt.Printf(
			"User permissions successfully copied from %s/%s to %s/%s\n",
			sourceRepo.Project.Key,
			sourceRepo.Slug,
			targetRepo.Project.Key,
			targetRepo.Slug,
		)
	}

	if command.flags.groupPermissions == true {
		err := client.CloneRepositoryGroupPermissions(sourceRepo.Project.Key, sourceRepo.Slug, targetRepo.Project.Key, targetRepo.Slug)
		if err != nil {
			return err
		}
		fmt.Printf(
			"Group permissions successfully copied from %s/%s to %s/%s\n",
			sourceRepo.Project.Key,
			sourceRepo.Slug,
			targetRepo.Project.Key,
			targetRepo.Slug,
		)
	}

	if command.flags.branchRestrictions == true {
		err := client.CloneRepositoryMasterBranchRestrictions(sourceRepo.Project.Key, sourceRepo.Slug, targetRepo.Project.Key, targetRepo.Slug)
		if err != nil {
			return err
		}
		fmt.Printf(
			"Branch restrictions successfully copied from %s/%s to %s/%s\n",
			sourceRepo.Project.Key,
			sourceRepo.Slug,
			targetRepo.Project.Key,
			targetRepo.Slug,
		)
	}

	if command.flags.pullRequestSettings == true {
		settings, err := client.GetPullRequestSettings(sourceRepo.Project.Key, sourceRepo.Slug)
		if err != nil {
			return err
		}

		err = client.SetPullRequestSettings(targetRepo.Project.Key, targetRepo.Slug, settings)
		if err != nil {
			return err
		}
		fmt.Printf(
			"Pull request settings successfully copied from %s/%s to %s/%s\n",
			sourceRepo.Project.Key,
			sourceRepo.Slug,
			targetRepo.Project.Key,
			targetRepo.Slug,
		)
	}

	return nil
}
