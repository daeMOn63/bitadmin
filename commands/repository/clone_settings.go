package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

type RepositoryCloneSettingsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *RepositoryCloneSettingsCommandFlags
}

type RepositoryCloneSettingsCommandFlags struct {
	sourceRepository   string
	targetRepository   string
	userPermissions    bool
	groupPermissions   bool
	branchRestrictions bool
}

func (command *RepositoryCloneSettingsCommand) GetCommand(fileCache *helper.FileCache) cli.Command {
	return cli.Command{
		Name:   "clone-settings",
		Usage:  "Clone various settings from a repository to another",
		Action: command.RepositoryCloneSettingsAction,
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
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, fileCache)
		},
	}
}

func (command *RepositoryCloneSettingsCommand) RepositoryCloneSettingsAction(context *cli.Context) error {
	client, err := command.Settings.GetApiClient()
	if err != nil {
		return err
	}

	fileCache := command.Settings.GetFileCache()

	var sourceRepo, targetRepo bitclient.Repository
	var sourceFound, targetFound bool

	for _, repo := range fileCache.Repositories {
		switch repo.Slug {
		case command.flags.sourceRepository:
			sourceFound = true
			sourceRepo = repo
		case command.flags.targetRepository:
			targetFound = true
			targetRepo = repo
		}
	}

	fmt.Println(sourceRepo.Name)
	fmt.Println(targetRepo.Name)

	if sourceFound == false {
		return fmt.Errorf("Cannot find repository %s from cache", command.flags.sourceRepository)
	}
	if targetFound == false {
		return fmt.Errorf("Cannot find repository %s from cache", command.flags.targetRepository)
	}

	if command.flags.userPermissions == true {
		err := client.CloneRepositoryUserPermissions(sourceRepo.Project.Key, sourceRepo.Slug, targetRepo.Project.Key, targetRepo.Slug)
		if err != nil {
			return err
		}
	}

	if command.flags.groupPermissions == true {
		err := client.CloneRepositoryGroupPermissions(sourceRepo.Project.Key, sourceRepo.Slug, targetRepo.Project.Key, targetRepo.Slug)
		if err != nil {
			return err
		}
	}

	if command.flags.branchRestrictions == true {
		err := client.CloneRepositoryMasterBranchRestrictions(sourceRepo.Project.Key, sourceRepo.Slug, targetRepo.Project.Key, targetRepo.Slug)
		if err != nil {
			return err
		}
	}

	return nil
}
