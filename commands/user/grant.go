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
	repositories cli.StringSlice
	usernames    cli.StringSlice
	permission   string
}

// GetCommand provide a ready to use cli.Command
func (command *GrantCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "grant",
		Usage:  "Grant users permission on repositories",
		Action: command.GrantAction,
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "repository",
				Usage: "The `<repository_slug>` the user will be added on",
				Value: &command.flags.repositories,
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

// GrantAction permit to give repository's permission to given users
func (command *GrantCommand) GrantAction(context *cli.Context) error {

	if len(command.flags.repositories) == 0 {
		return fmt.Errorf("flag --repository is required")
	}

	if len(command.flags.usernames) == 0 {
		return fmt.Errorf("At least one --username is required")
	}

	if len(command.flags.permission) == 0 {
		return fmt.Errorf("flag --permission is required")
	}

	fileCache := command.Settings.GetFileCache()

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	for _, repositorySlug := range command.flags.repositories {

		repo, err := fileCache.SearchRepositorySlug(repositorySlug)

		if err != nil {
			return err
		}

		for _, username := range command.flags.usernames {
			params := bitclient.SetRepositoryUserPermissionRequest{
				Username:   username,
				Permission: command.flags.permission,
			}

			err := client.SetRepositoryUserPermission(repo.Project.Key, repositorySlug, params)

			if err != nil {
				fmt.Printf("[KO] rep%s - %s\n", username, err)
				return fmt.Errorf("repo %s, user %s, permission %s - reason: %s", repositorySlug, username, command.flags.permission, err)
			}

			fmt.Printf("[OK] repo %s, user %s, permission %s\n", repositorySlug, username, command.flags.permission)
		}
	}

	fmt.Printf("Done granting user permissions\n")

	return nil
}
