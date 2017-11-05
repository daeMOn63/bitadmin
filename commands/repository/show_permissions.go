// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// ShowPermissionsCommand define base struct for ShowPermissions actions
type ShowPermissionsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *ShowPermissionsFlags
}

// ShowPermissionsFlags define flags required by the ShowPermissionsAction
type ShowPermissionsFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *ShowPermissionsCommand) GetCommand(fileCache *helper.FileCache) cli.Command {
	return cli.Command{
		Name:   "show-permission",
		Usage:  "Show permissions on given repository",
		Action: command.ShowPermissionsAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<rproject>` of the repository",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository>` which to dump permissions from",
				Destination: &command.flags.repository,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, fileCache)
		},
	}
}

// ShowPermissionsAction display the current user / group permissions on given repository
func (command *ShowPermissionsCommand) ShowPermissionsAction(context *cli.Context) error {

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	userResponse, err := client.GetRepositoryUserPermission(
		command.flags.project,
		command.flags.repository,
		bitclient.GetRepositoryUserPermissionRequest{},
	)

	if err != nil {
		return err
	}

	groupResponse, err := client.GetRepositoryGroupPermission(
		command.flags.project,
		command.flags.repository,
		bitclient.GetRepositoryGroupPermissionRequest{},
	)

	if err != nil {
		return err
	}

	if len(userResponse.Values) <= 0 {
		fmt.Printf("No user permissions found on repository %s\n", command.flags.repository)
	}

	for _, userPermission := range userResponse.Values {
		fmt.Printf("user %s - %s\n", userPermission.User.Slug, userPermission.Permission)
	}

	if len(groupResponse.Values) <= 0 {
		fmt.Printf("No group permissions found on repository %s\n", command.flags.repository)
	}

	for _, groupPermission := range groupResponse.Values {
		fmt.Printf("group %s - %s\n", groupPermission.Group.Name, groupPermission.Permission)
	}

	return nil
}
