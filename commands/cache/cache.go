package cache

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

type CacheCommand struct {
	Settings *settings.BitAdminSettings
}

func (cc *CacheCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:  "cache",
		Usage: "Caching data for faster operation",
		Subcommands: []cli.Command{
			{
				Name:   "clear",
				Usage:  "Clear all cached data",
				Action: cc.ClearCacheAction,
			},
			{
				Name:   "warmup",
				Usage:  "Fetch data and cache them",
				Action: cc.WarmupCacheAction,
			},
			{
				Name:   "dump",
				Usage:  "Print current cache content",
				Action: cc.DumpCacheAction,
			},
		},
	}
}

func (command *CacheCommand) ClearCacheAction(context *cli.Context) error {
	fmt.Println("clearing cache... ")
	return command.Settings.GetFileCache().Clear()
}

func (command *CacheCommand) WarmupCacheAction(context *cli.Context) error {
	fmt.Println("warming up cache... ")
	client, err := command.Settings.GetApiClient()
	if err != nil {
		return err
	}

	cache := command.Settings.GetFileCache()
	if err != nil {
		return err
	}

	maxPagedRequest := bitclient.PagedRequest{
		Limit: 10000,
		Start: 0,
	}

	userResponse, err := client.GetUsers(maxPagedRequest)
	if err != nil {
		return err
	}
	cache.Users = userResponse.Values

	projectResponse, err := client.GetProjects(maxPagedRequest)
	if err != nil {
		return err
	}
	cache.Projects = projectResponse.Values

	for _, project := range cache.Projects {

		repositoryResponse, err := client.GetRepositories(project.Key, maxPagedRequest)

		if err != nil {
			return err
		}

		cache.Repositories = append(cache.Repositories, repositoryResponse.Values...)
	}

	return cache.Save()
}

func (command *CacheCommand) DumpCacheAction(context *cli.Context) error {
	fmt.Println(command.Settings.GetFileCache())
	return nil
}
