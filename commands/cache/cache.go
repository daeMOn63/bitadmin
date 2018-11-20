// Package cache provide actions for loading / clearing / dumping the users, repositories, groups, projects from Bitbucket
// It aims to provide fluid autocompletion and avoid hitting the API while searching for specific entities.
package cache

import (
	"fmt"

	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// Command define the root command holding cache actions
type Command struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *Command) GetCommand() cli.Command {
	return cli.Command{
		Name:  "cache",
		Usage: "Caching data for faster operation and autocompletion.",
		Subcommands: []cli.Command{
			{
				Name:   "clear",
				Usage:  "Clear all cached data",
				Action: command.ClearCacheAction,
			},
			{
				Name:   "warmup",
				Usage:  "Fetch data and cache them",
				Action: command.WarmupCacheAction,
			},
			{
				Name:   "dump",
				Usage:  "Print current cache content",
				Action: command.DumpCacheAction,
			},
		},
	}
}

// ClearCacheAction wipe out the cache
func (command *Command) ClearCacheAction(context *cli.Context) error {
	fmt.Println("Clearing cache... ")
	return command.Settings.GetFileCache().Clear()
}

// WarmupCacheAction load all entities and save them into a file
func (command *Command) WarmupCacheAction(context *cli.Context) error {
	fmt.Println("Warming up cache...")
	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	cache := command.Settings.GetFileCache()

	var limit uint
	var offset uint
	var isLastPage bool

	cache.Clear()

	fmt.Printf("Loading users...")
	limit = 1000
	offset = uint(0)
	isLastPage = false

	for !isLastPage {
		userResponse, err := client.GetUsers(bitclient.PagedRequest{
			Limit: limit,
			Start: offset,
		})
		if err != nil {
			return err
		}

		cache.Users = append(cache.Users, userResponse.Values...)

		isLastPage = userResponse.IsLastPage

		offset += limit

	}
	fmt.Println("done")
	fmt.Printf("Cached %d users\n", len(cache.Users))

	fmt.Printf("Loading projects...")
	fmt.Printf("Loading users...")
	limit = 1000
	offset = uint(0)
	isLastPage = false
	for !isLastPage {
		projectResponse, err := client.GetProjects(bitclient.PagedRequest{
			Limit: limit,
			Start: offset,
		})
		if err != nil {
			return err
		}
		cache.Projects = append(cache.Projects, projectResponse.Values...)

		isLastPage = projectResponse.IsLastPage

		offset += limit
	}
	fmt.Println("done")
	fmt.Printf("Cached %d projects\n", len(cache.Projects))

	fmt.Printf("Loading repositories...")
	for _, project := range cache.Projects {

		limit = 1000
		offset = uint(0)
		isLastPage = false

		for !isLastPage {
			repositoryResponse, err := client.GetRepositories(project.Key, bitclient.PagedRequest{
				Limit: limit,
				Start: offset,
			})

			if err != nil {
				return err
			}

			cache.Repositories = append(cache.Repositories, repositoryResponse.Values...)

			isLastPage = repositoryResponse.IsLastPage

			offset += limit
		}
	}
	fmt.Println("done")
	fmt.Printf("Cached %d repositories\n", len(cache.Repositories))

	fmt.Println("\nCache warmup completed")
	return cache.Save()
}

// DumpCacheAction print on stdout the content of the cache
func (command *Command) DumpCacheAction(context *cli.Context) error {
	fmt.Println(command.Settings.GetFileCache())
	return nil
}
