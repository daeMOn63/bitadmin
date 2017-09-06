package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

type RepositoryUserCommand struct {
	Settings *settings.BitAdminSettings
	flags    *RepositoryUserCommandFlags
}

type RepositoryUserCommandFlags struct {
	repository string
	usernames  cli.StringSlice
	permission string
}

func (command *RepositoryUserCommand) GetCommand(fileCache *helper.FileCache) cli.Command {
	return cli.Command{
		Name:  "user",
		Usage: "User operations on the repository",
		Subcommands: []cli.Command{
			{
				Name:   "add",
				Usage:  "Add a user on the repository",
				Action: command.AddUserAction,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:        "repository",
						Usage:       "The `<repository_name>` the user will be added on",
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
					helper.AutoComplete(c, fileCache)
				},
			},
		},
	}
}

func (command *RepositoryUserCommand) AddUserAction(context *cli.Context) error {

	if len(command.flags.repository) == 0 {
		return fmt.Errorf("flag --repository is required.")
	}

	if len(command.flags.usernames) == 0 {
		return fmt.Errorf("flag --username is required.")
	}

	if len(command.flags.permission) == 0 {
		return fmt.Errorf("flag --permission is required.")
	}

	fileCache := command.Settings.GetFileCache()
	fileCache.Load()

	repo, err := fileCache.SearchRepositorySlug(command.flags.repository)

	if err != nil {
		return err
	}

	projectKey := repo.Project.Key

	client, err := command.Settings.GetApiClient()
	if err != nil {
		return err
	}

	fmt.Printf("\nGranting permission %s on repository %s (%s)\n\n", command.flags.permission, repo.Name, repo.Links["self"][0]["href"])

	failedCount := 0
	for _, username := range command.flags.usernames {
		err := client.SetRepositoryUserPermission(projectKey, command.flags.repository, username, command.flags.permission)

		if err != nil {
			fmt.Printf("[KO] %s - %s\n", username, err)
			failedCount += 1
		} else {
			fmt.Printf("[OK] %s\n", username)
		}
	}

	if failedCount > 0 {
		return fmt.Errorf("%d user(s) have not been added properly.\n", failedCount)
	} else {
		fmt.Println("\nSuccess: All users have been granted permissions.\n")
	}

	return nil
}
