// Package repository hold actions on the Bitbucket repositories
package user

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// UnsetPermissionsCommand define base struct for Move actions
type UnsetPermissionsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *UnsetPermissionsCommandFlags
}

// UnsetPermissionsCommandFlags hold flag values for the UnsetPermissionsCommand
type UnsetPermissionsCommandFlags struct {
	project    string
	repository string
	usernames  cli.StringSlice
}

// GetCommand provide a ready to use cli.Command
func (command *UnsetPermissionsCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "unset-permissions",
		Usage:  "Unset user permissions on given repository",
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
				Name:  "username",
				Usage: "The `<username>` to unset permissions for. Can be repeated multiple times",
				Value: &command.flags.usernames,
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

	if len(command.flags.usernames) == 0 {
		return fmt.Errorf("At least one --username is required")
	}

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}
	for _, username := range command.flags.usernames {
		params := bitclient.UnsetRepositoryUserPermissionRequest{
			Username: username,
		}

		err := client.UnsetRepositoryUserPermission(command.flags.project, command.flags.repository, params)

		if err != nil {
			return fmt.Errorf(
				"Cannot unset permissions for repo %s/%s, user %s - reason: %s",
				command.flags.project,
				command.flags.repository,
				username,
				err,
			)
		}

		fmt.Printf(
			"[OK] Permissions removed on repo %s/%s, user %s\n",
			command.flags.project,
			command.flags.repository,
			username,
		)
	}

	return nil
}
