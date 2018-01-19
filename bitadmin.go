package main

import (
	"os"
	"sort"

	"github.com/daeMOn63/bitadmin/commands/cache"
	"github.com/daeMOn63/bitadmin/commands/group"
	"github.com/daeMOn63/bitadmin/commands/hooks"
	"github.com/daeMOn63/bitadmin/commands/repository"
	"github.com/daeMOn63/bitadmin/commands/user"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/fatih/color"
	"github.com/urfave/cli"
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

	cacheCommand := &cache.Command{
		Settings: globalSettings,
	}

	repositoryCommand := &repository.Command{
		Settings: globalSettings,
	}

	userCommand := &user.Command{
		Settings: globalSettings,
	}

	groupCommand := &group.Command{
		Settings: globalSettings,
	}

	hooksCommand := &hooks.Command{
		Settings: globalSettings,
	}

	app.Commands = []cli.Command{
		cacheCommand.GetCommand(),
		repositoryCommand.GetCommand(),
		userCommand.GetCommand(),
		groupCommand.GetCommand(),
		hooksCommand.GetCommand(),
	}

	app.BashComplete = helper.AppAutoComplete

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)

	if err != nil {
		color.New(color.FgRed).Fprintf(cli.ErrWriter, "\nError: %s\n", err.Error())
		os.Exit(1)
	}
}
