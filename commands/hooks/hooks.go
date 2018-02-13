// Package hooks hold actions on the Bitbucket hooks
package hooks

import (
	"errors"
	"fmt"

	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// Command define base struct for repository subcommands and actions
type Command struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *Command) GetCommand() cli.Command {

	listHookCommand := ListHooksCommand{
		Settings: command.Settings,
		flags:    &ListHooksCommandFlags{},
	}

	eolHookCommand := EolHookCommand{
		Settings: command.Settings,
	}

	pubHookCommand := PubHookCommand{
		Settings: command.Settings,
	}

	yaccHookCommand := YaccHookCommand{
		Settings: command.Settings,
	}

	rfpHookCommand := RfpHookCommand{
		Settings: command.Settings,
	}

	return cli.Command{
		Name:  "hooks",
		Usage: "Hooks operations",
		Subcommands: []cli.Command{
			listHookCommand.GetCommand(),
			eolHookCommand.GetCommand(),
			pubHookCommand.GetCommand(),
			yaccHookCommand.GetCommand(),
			rfpHookCommand.GetCommand(),
		},
	}
}

// ListHooksCommand define command to retrieve hooks on a given repository
type ListHooksCommand struct {
	Settings *settings.BitAdminSettings
	flags    *ListHooksCommandFlags
}

// ListHooksCommandFlags define the flags for the ListHooksCommand
type ListHooksCommandFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *ListHooksCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "list",
		Usage:  "List hooks on repository",
		Action: command.ListHooksAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<project_key>` of the repository",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository_name>` to list hooks for",
				Destination: &command.flags.repository,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// ListHooksAction contains logic to list hooks on given repository
func (command *ListHooksCommand) ListHooksAction(context *cli.Context) error {
	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	if len(command.flags.project) <= 0 {
		return errors.New("--project flag is required")
	}
	if len(command.flags.repository) <= 0 {
		return errors.New("--repository flag is required")
	}

	response, err := client.GetHooks(command.flags.project, command.flags.repository, bitclient.GetHooksRequest{})

	if err != nil {
		return err
	}

	for _, hook := range response.Values {

		status := "DISABLED"
		if hook.Enabled == true {
			status = "ENABLED "
		}

		fmt.Printf("[%s] %s (%s)\n", status, hook.Details.Name, hook.Details.Key)
	}

	return nil
}
