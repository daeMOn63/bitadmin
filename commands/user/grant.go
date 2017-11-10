// Package user hold the actions on the Bitbucket users
package user

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// GrantCommand define base struct for user grant action
type GrantCommand struct {
	Settings *settings.BitAdminSettings
	flags    *GrantCommandFlags
}

// GrantCommandFlags hold the flag values of the use grant action
type GrantCommandFlags struct {
	project    string
	repository string
	usernames  cli.StringSlice
	permission string
}

// GetCommand provide a ready to use cli.Command
func (command *GrantCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "grant",
		Usage:  "Grant users permission on repositories",
		Action: command.GrantAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<rproject>` of the repository",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository_slug>` the user will be added on",
				Destination: &command.flags.repository,
			},
			cli.StringSliceFlag{
				Name:  "username",
				Usage: "The `<username>` to be added on the repository. Can be repeated multiple times",
				Value: &command.flags.usernames,
			},
			cli.StringFlag{
				Name:        "permission",
				Usage:       "The `<permission>` level the user will have (one of REPO_READ, REPO_WRITE, REPO_ADMIN)",
				Destination: &command.flags.permission,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// GrantAction permit to grant permission on repository to given users
func (command *GrantCommand) GrantAction(context *cli.Context) error {

	if len(command.flags.project) == 0 {
		return fmt.Errorf("flag --project is required")
	}

	if len(command.flags.repository) == 0 {
		return fmt.Errorf("flag --repository is required")
	}

	if len(command.flags.usernames) == 0 {
		return fmt.Errorf("At least one --username is required")
	}

	if len(command.flags.permission) == 0 {
		return fmt.Errorf("flag --permission is required")
	}

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	for _, username := range command.flags.usernames {
		params := bitclient.SetRepositoryUserPermissionRequest{
			Username:   username,
			Permission: command.flags.permission,
		}

		err := client.SetRepositoryUserPermission(command.flags.project, command.flags.repository, params)

		if err != nil {
			return fmt.Errorf(
				"repo %s/%s, user %s, permission %s - reason: %s",
				command.flags.project,
				command.flags.repository,
				username,
				command.flags.permission,
				err,
			)
		}

		fmt.Printf(
			"[OK] repo %s/%s, user %s, permission %s\n",
			command.flags.project,
			command.flags.repository,
			username,
			command.flags.permission,
		)
	}

	return nil
}
