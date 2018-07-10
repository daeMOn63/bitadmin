// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// MoveCommand define base struct for Move actions
type MoveCommand struct {
	Settings *settings.BitAdminSettings
	flags    *MoveCommandFlags
}

// MoveCommandFlags hold flag values for the MoveCommand
type MoveCommandFlags struct {
	project          string
	repository       string
	targetProject    string
	targetRepository string
}

// GetCommand provide a ready to use cli.Command
func (command *MoveCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "move",
		Usage:  "Move / rename a repository",
		Action: command.MoveRepositoryAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<prroject>` of the repository",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository>` to move",
				Destination: &command.flags.repository,
			},
			cli.StringFlag{
				Name:        "targetProject",
				Usage:       "The target `<project>` where the repository will get moved",
				Destination: &command.flags.targetProject,
			},
			cli.StringFlag{
				Name:        "targetRepository",
				Usage:       "The target `<repository>` where the repository will get renamed",
				Destination: &command.flags.targetRepository,
			},
		},
		BashComplete: func(c *cli.Context) {
			// TODO: improve autocomplete. Here it's not perfect as we want targetProject / targetRepo
			// flag to list everything, ignoring what can already be set in project / repository flags.
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// MoveRepositoryAction use flag values to create a new repository
func (command *MoveCommand) MoveRepositoryAction(context *cli.Context) error {
	if len(command.flags.project) == 0 {
		return fmt.Errorf("flag --project is required")
	}

	if len(command.flags.repository) == 0 {
		return fmt.Errorf("flag --repository is required")
	}

	if len(command.flags.targetProject) == 0 && len(command.flags.targetRepository) == 0 {
		return fmt.Errorf("at least one flag --targetProject or --targetRepository is required")
	}

	if len(command.flags.targetRepository) == 0 {
		command.flags.targetRepository = command.flags.repository
	}

	if len(command.flags.targetProject) == 0 {
		command.flags.targetProject = command.flags.project
	}

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	params := bitclient.UpdateRepositoryRequest{
		Name:    command.flags.targetRepository,
		Project: bitclient.Project{Key: command.flags.targetProject},
	}

	client.UpdateRepository(command.flags.project, command.flags.repository, params)
	if err != nil {
		return err
	}

	fmt.Printf("[OK]%s/%s moved to %s/%s\n",
		command.flags.project,
		command.flags.repository,
		command.flags.targetProject,
		command.flags.targetRepository,
	)

	return nil
}
