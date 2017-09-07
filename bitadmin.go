package main

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/commands/cache"
	"github.com/daeMOn63/bitadmin/commands/repository"
	"github.com/daeMOn63/bitadmin/commands/user"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
	"os"
	"sort"
)

func main() {

	globalSettings := settings.NewSettings()
	app := cli.NewApp()

	app.EnableBashCompletion = true

	app.Name = "bitadmin"
	app.Usage = "Bitbucket cli administration tool"
	app.Version = "0.0.1"
	app.Author = "Flavien Binet"
	app.Email = "https://github.com/daeMOn63/bitadmin"
	app.Flags = globalSettings.GetFlags()

	cacheCommand := &cache.CacheCommand{
		Settings: globalSettings,
	}

	repositoryCommand := &repository.RepositoryCommand{
		Settings: globalSettings,
	}

	userCommand := &user.UserCommand{
		Settings: globalSettings,
	}

	app.Commands = []cli.Command{
		cacheCommand.GetCommand(),
		repositoryCommand.GetCommand(),
		userCommand.GetCommand(),
	}

	app.BashComplete = helper.AppAutoComplete

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)

	if err != nil {
		fmt.Fprintf(cli.ErrWriter, "\nError: %s\n", err.Error())
		os.Exit(1)
	}
}
