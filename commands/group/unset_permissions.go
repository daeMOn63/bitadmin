// Package repository hold actions on the Bitbucket repositories
package group

import (
	"fmt"

	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// UnsetPermissionsCommand define base struct for Unset permissions actions
type UnsetPermissionsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *UnsetPermissionsCommandFlags
}

// UnsetPermissionsCommandFlags hold flag values for the UnsetPermissionsCommand
type UnsetPermissionsCommandFlags struct {
	project    string
	repository string
	groups     cli.StringSlice
}

// GetCommand provide a ready to use cli.Command
func (command *UnsetPermissionsCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "unset-permissions",
		Usage:  "Unset groups permissions on given repository",
		Action: command.UnsetPermissionsAction,
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
			cli.StringSliceFlag{
				Name:  "name",
				Usage: "The `<name>` of the group to unset permissions for. Can be repeated multiple times",
				Value: &command.flags.groups,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// UnsetPermissionsAction unset permissions of given user(s) on given repository
func (command *UnsetPermissionsCommand) UnsetPermissionsAction(context *cli.Context) error {

	if len(command.flags.project) == 0 {
		return fmt.Errorf("flag --project is required")
	}

	if len(command.flags.repository) == 0 {
		return fmt.Errorf("flag --repository is required")
	}

	if len(command.flags.groups) == 0 {
		return fmt.Errorf("At least one --name is required")
	}

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}
	for _, group := range command.flags.groups {
		params := bitclient.UnsetRepositoryGroupPermissionRequest{
			Name: group,
		}

		err := client.UnsetRepositoryGroupPermission(command.flags.project, command.flags.repository, params)

		if err != nil {
			return fmt.Errorf(
				"Cannot unset permissions for repo %s/%s, group %s - reason: %s",
				command.flags.project,
				command.flags.repository,
				group,
				err,
			)
		}

		fmt.Printf(
			"[OK] Permissions removed on repo %s/%s, group %s\n",
			command.flags.project,
			command.flags.repository,
			group,
		)
	}

	return nil
}
