// Package hooks hold actions on the Bitbucket hooks
package hooks

import (
	"errors"
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

const rfpHookKey = "com.atlassian.bitbucket.server.bitbucket-bundled-hooks:force-push-hook"

// RfpHookCommand define the command for the Reject Force Push hook
type RfpHookCommand struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *RfpHookCommand) GetCommand() cli.Command {

	pubEnableCommand := RfpHookEnableCommand{
		Settings: command.Settings,
		flags:    &RfpHookEnableCommandFlags{},
	}

	pubDisableCommand := RfpHookDisableCommand{
		Settings: command.Settings,
		flags:    &RfpHookDisableCommandFlags{},
	}

	return cli.Command{
		Name:  "reject-force-push",
		Usage: "Reject Force Push hook operations",
		Subcommands: []cli.Command{
			pubEnableCommand.GetCommand(),
			pubDisableCommand.GetCommand(),
		},
	}
}

// RfpHookEnableCommand define the command to enable the PUB hook
type RfpHookEnableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *RfpHookEnableCommandFlags
}

// RfpHookEnableCommandFlags define the flags for the RfpHookEnableCommand
type RfpHookEnableCommandFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *RfpHookEnableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "enable",
		Usage:  "Enable Reject Force Push hook",
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
func (command *RfpHookEnableCommand) EnableAction(context *cli.Context) error {
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
		rfpHookKey,
		nil,
	)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Enabled and configured Reject Force Push hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}

// RfpHookDisableCommand define the command to disable the YACC hook
type RfpHookDisableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *RfpHookDisableCommandFlags
}

// RfpHookDisableCommandFlags define the flags of the RfpHookDisableCommand
type RfpHookDisableCommandFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *RfpHookDisableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "disable",
		Usage:  "Disable Reject Force Push hook",
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
func (command *RfpHookDisableCommand) DisableAction(context *cli.Context) error {
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

	err = client.DisableHook(command.flags.project, command.flags.repository, rfpHookKey)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Disabled Reject Force Push hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}
