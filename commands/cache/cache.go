package cache

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/settings"
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

	command.Settings.GetFileCache().Users, err = client.GetUsers()
	if err != nil {
		return err
	}

	command.Settings.GetFileCache().Projects, err = client.GetProjects()
	if err != nil {
		return err
	}

	for _, project := range command.Settings.GetFileCache().Projects {

		repositories, err := client.GetRepositories(project.Key)

		if err != nil {
			return err
		}

		command.Settings.GetFileCache().Repositories = append(command.Settings.GetFileCache().Repositories, repositories...)
	}

	return command.Settings.GetFileCache().Save()
}

func (command *CacheCommand) DumpCacheAction(context *cli.Context) error {
	err := command.Settings.GetFileCache().Load()
	fmt.Println(command.Settings.GetFileCache())
	return err
}
