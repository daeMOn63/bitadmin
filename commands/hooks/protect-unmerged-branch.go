// Package hooks hold actions on the Bitbucket hooks
package hooks

import (
	"errors"
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

const pubHookKey = "com.atlassian.stash.plugin.stash-protect-unmerged-branch-hook:protect-unmerged-branch-hook"

// PubHookCommand define the command for the Protect Unmerged Branch hook
type PubHookCommand struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *PubHookCommand) GetCommand() cli.Command {

	pubEnableCommand := PubHookEnableCommand{
		Settings: command.Settings,
		flags:    &PubHookEnableCommandFlags{},
	}

	pubDisableCommand := PubHookDisableCommand{
		Settings: command.Settings,
		flags:    &PubHookDisableCommandFlags{},
	}

	return cli.Command{
		Name:  "protect-unmerged-branch",
		Usage: "Protect Unmerged Branch hook operations",
		Subcommands: []cli.Command{
			pubEnableCommand.GetCommand(),
			pubDisableCommand.GetCommand(),
		},
	}
}

// PubHookEnableCommand define the command to enable the PUB hook
type PubHookEnableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *PubHookEnableCommandFlags
}

// PubHookEnableCommandFlags define the flags for the PubHookEnableCommand
type PubHookEnableCommandFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *PubHookEnableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "enable",
		Usage:  "Enable Protect Unmerged Branch hook and set its configuration",
		Action: command.EnableAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<project_key>` containing the repository to enable the hook on",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository_name>` to enable the hook on",
				Destination: &command.flags.repository,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// EnableAction contains logic to turn on the hook and set its configuration
func (command *PubHookEnableCommand) EnableAction(context *cli.Context) error {
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

	err = client.EnableHook(
		command.flags.project,
		command.flags.repository,
		pubHookKey,
		nil,
	)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Enabled and configured Protect Unmerged Branch hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}

// PubHookDisableCommand define the command to disable the YACC hook
type PubHookDisableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *PubHookDisableCommandFlags
}

// PubHookDisableCommandFlags define the flags of the PubHookDisableCommand
type PubHookDisableCommandFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *PubHookDisableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "disable",
		Usage:  "Disable Protect Unmerged Branch hook",
		Action: command.DisableAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<project_key>` where the repository will be created",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository_name>` to create",
				Destination: &command.flags.repository,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// DisableAction contains logic to turn on the hook and set its configuration
func (command *PubHookDisableCommand) DisableAction(context *cli.Context) error {
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

	err = client.DisableHook(command.flags.project, command.flags.repository, pubHookKey)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Disabled Protect Unmerged Branch hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}
