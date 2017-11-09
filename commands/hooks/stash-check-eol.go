// Package hooks hold actions on the Bitbucket hooks
package hooks

import (
	"errors"
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

const eolHookKey = "com.pbaranchikov.stash-eol-check:stash-check-eol-hook"

// EolHookCommand define the command for the Stash Check EOL hook
type EolHookCommand struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *EolHookCommand) GetCommand() cli.Command {

	eolEnableCommand := EolHookEnableCommand{
		Settings: command.Settings,
		flags:    &EolHookEnableCommandFlags{},
	}

	eolDisableCommand := EolHookDisableCommand{
		Settings: command.Settings,
		flags:    &EolHookDisableCommandFlags{},
	}

	return cli.Command{
		Name:  "stash-eol-check",
		Usage: "Stash EOL Check hook operations",
		Subcommands: []cli.Command{
			eolEnableCommand.GetCommand(),
			eolDisableCommand.GetCommand(),
		},
	}
}

// EolHookEnableCommand define the command to enable the EOL hook
type EolHookEnableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *EolHookEnableCommandFlags
}

// EolSettings define the settings of the EOL hook
type EolSettings struct {
	ExcludeFiles      string `json:"excludeFiles,omitempty"`
	AllowInheritedEol bool   `json:"allowInheritedEol,omitempty"`
}

// EolHookEnableCommandFlags define the flags for the EolHookEnableCommand
type EolHookEnableCommandFlags struct {
	project    string
	repository string
	settings   EolSettings
}

// GetCommand provide a ready to use cli.Command
func (command *EolHookEnableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "enable",
		Usage:  "Enable Stash Check EOL hook and set its configuration",
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
			cli.StringFlag{
				Name:        "excludeFiles",
				Usage:       "`<files>` that should not be checked for EOL - comma separated list of regexps",
				Destination: &command.flags.settings.ExcludeFiles,
			},
			cli.BoolFlag{
				Name:        "allowInheritedEol",
				Usage:       "allow commit of wrong EOL-style for files, that are already committed with wrong EOL-style",
				Destination: &command.flags.settings.AllowInheritedEol,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// EnableAction contains logic to turn on the hook and set its configuration
func (command *EolHookEnableCommand) EnableAction(context *cli.Context) error {
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
		eolHookKey,
		command.flags.settings,
	)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Enabled and configured eol hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}

// EolHookDisableCommand define the command to disable the YACC hook
type EolHookDisableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *EolHookDisableCommandFlags
}

// EolHookDisableCommandFlags define the flags of the EolHookDisableCommand
type EolHookDisableCommandFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *EolHookDisableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "disable",
		Usage:  "Disable Stash Check EOL hook",
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
func (command *EolHookDisableCommand) DisableAction(context *cli.Context) error {
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

	err = client.DisableHook(command.flags.project, command.flags.repository, eolHookKey)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Disabled eol hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}
