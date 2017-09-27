package user

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

type UserGrantCommand struct {
	Settings *settings.BitAdminSettings
	flags    *UserGrantCommandFlags
}

type UserGrantCommandFlags struct {
	repositories cli.StringSlice
	usernames    cli.StringSlice
	permission   string
}

func (command *UserGrantCommand) GetCommand() cli.Command {
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

func (command *UserGrantCommand) GrantAction(context *cli.Context) error {

	if len(command.flags.repositories) == 0 {
		return fmt.Errorf("flag --repository is required.")
	}

	if len(command.flags.usernames) == 0 {
		return fmt.Errorf("flag --username is required.")
	}

	if len(command.flags.permission) == 0 {
		return fmt.Errorf("flag --permission is required.")
	}

	fileCache := command.Settings.GetFileCache()

	client, err := command.Settings.GetApiClient()
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
				return fmt.Errorf("repo %s, user %s, permission %s - reason: %s\n", repositorySlug, username, command.flags.permission, err)
			} else {
				fmt.Printf("[OK] repo %s, user %s, permission %s\n", repositorySlug, username, command.flags.permission)
			}
		}
	}

	fmt.Printf("Done granting permissions\n")

	return nil
}
