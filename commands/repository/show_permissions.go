package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

type RepositoryShowPermissionsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *RepositoryShowPermissionsFlags
}

type RepositoryShowPermissionsFlags struct {
	repositorySlug string
}

func (command *RepositoryShowPermissionsCommand) GetCommand(fileCache *helper.FileCache) cli.Command {
	return cli.Command{
		Name:   "show-permission",
		Usage:  "Show permissions on given repository",
		Action: command.ShowRepositoryPermissionsAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository>` which to dump permissions from",
				Destination: &command.flags.repositorySlug,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, fileCache)
		},
	}
}

func (command *RepositoryShowPermissionsCommand) ShowRepositoryPermissionsAction(context *cli.Context) error {

	fileCache := command.Settings.GetFileCache()

	client, err := command.Settings.GetApiClient()
	if err != nil {
		return err
	}

	var repository *bitclient.Repository
	for _, repo := range fileCache.Repositories {
		if repo.Slug == command.flags.repositorySlug {
			repository = &repo
			break
		}
	}

	if repository == nil {
		return fmt.Errorf("Cannot retreive repository %s", command.flags.repositorySlug)
	}

	userResponse, err := client.GetRepositoryUserPermission(repository.Project.Key, repository.Slug, bitclient.GetRepositoryUserPermissionRequest{})
	if err != nil {
		return err
	}

	groupResponse, err := client.GetRepositoryGroupPermission(repository.Project.Key, repository.Slug, bitclient.GetRepositoryGroupPermissionRequest{})
	if err != nil {
		return err
	}

	if len(userResponse.Values) <= 0 {
		fmt.Printf("No user permissions found on repository %s\n", command.flags.repositorySlug)
	}

	for _, userPermission := range userResponse.Values {
		fmt.Printf("user %s - %s\n", userPermission.User.Slug, userPermission.Permission)
	}

	if len(groupResponse.Values) <= 0 {
		fmt.Printf("No group permissions found on repository %s\n", command.flags.repositorySlug)
	}

	for _, groupPermission := range groupResponse.Values {
		fmt.Printf("group %s - %s\n", groupPermission.Group.Name, groupPermission.Permission)
	}

	return nil
}
