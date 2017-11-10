// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"errors"
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// BranchingModelCommand define base struct for BranchingModel actions
type BranchingModelCommand struct {
	Settings *settings.BitAdminSettings
	flags    *BranchingModelCommandFlags
}

// BranchingModelCommandFlags define flags required by the ShowPermissionsAction
type BranchingModelCommandFlags struct {
	project          string
	repository       string
	enableBugfix     bool
	enableFeature    bool
	enableHotfix     bool
	enableRelease    bool
	prefixBugfix     string
	prefixFeature    string
	prefixHotfix     string
	prefixRelease    string
	productionRefID  string
	developmentRefID string
}

// GetCommand provide a ready to use cli.Command
func (command *BranchingModelCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "set-branching-model",
		Usage:  "Set branching model options on given repository",
		Action: command.SetBranchingModelAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<rproject>` of the repository",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "repository",
				Usage:       "The `<repository>` which to dump permissions from",
				Destination: &command.flags.repository,
			},
			cli.BoolFlag{
				Name:        "enable-bugfix",
				Usage:       "Turn on the bugfix branch model",
				Destination: &command.flags.enableBugfix,
			},
			cli.BoolFlag{
				Name:        "enable-feature",
				Usage:       "Turn on the feature branch model",
				Destination: &command.flags.enableFeature,
			},
			cli.BoolFlag{
				Name:        "enable-hotfix",
				Usage:       "Turn on the hotfix branch model",
				Destination: &command.flags.enableHotfix,
			},
			cli.BoolFlag{
				Name:        "enable-release",
				Usage:       "Turn on the release branch model",
				Destination: &command.flags.enableRelease,
			},
			cli.StringFlag{
				Name:        "prefix-bugfix",
				Usage:       "Override default `<prefix>` for bugfix branch model",
				Destination: &command.flags.prefixBugfix,
			},
			cli.StringFlag{
				Name:        "prefix-feature",
				Usage:       "Override default `<prefix>` for feature branch model",
				Destination: &command.flags.prefixFeature,
			},
			cli.StringFlag{
				Name:        "prefix-hotfix",
				Usage:       "Override default `<prefix>` for hotfix branch model",
				Destination: &command.flags.prefixHotfix,
			},
			cli.StringFlag{
				Name:        "prefix-release",
				Usage:       "Override default `<prefix>` for release branch model",
				Destination: &command.flags.prefixRelease,
			},
			cli.StringFlag{
				Name:        "production-refid",
				Usage:       "Set the default `<production>` branch (ie: refs/heads/master)",
				Destination: &command.flags.productionRefID,
			},
			cli.StringFlag{
				Name:        "development-refid",
				Usage:       "Set the default `<development>` branch (ie: refs/heads/master)",
				Destination: &command.flags.developmentRefID,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// SetBranchingModelAction use flag values to set the branching model options on given repository
func (command *BranchingModelCommand) SetBranchingModelAction(context *cli.Context) error {

	if len(command.flags.project) <= 0 {
		return errors.New("--project flag is required")
	}
	if len(command.flags.repository) <= 0 {
		return errors.New("--repository flag is required")
	}

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	branchingModel, err := client.GetBranchingModel(
		command.flags.project,
		command.flags.repository,
	)

	if err != nil {
		return err
	}

	if len(command.flags.productionRefID) >= 0 {
		branchingModel.Production.RefId = command.flags.productionRefID
	}
	if len(command.flags.developmentRefID) >= 0 {
		branchingModel.Development.RefId = command.flags.developmentRefID
	}

	for i, t := range branchingModel.Types {
		switch t.Id {
		case "BUGFIX":
			branchingModel.Types[i].Enabled = command.flags.enableBugfix
			if len(command.flags.prefixBugfix) > 0 {
				branchingModel.Types[i].Prefix = command.flags.prefixBugfix
			}
		case "FEATURE":
			branchingModel.Types[i].Enabled = command.flags.enableFeature
			if len(command.flags.prefixFeature) > 0 {
				branchingModel.Types[i].Prefix = command.flags.prefixFeature
			}
		case "HOTFIX":
			branchingModel.Types[i].Enabled = command.flags.enableHotfix
			if len(command.flags.prefixHotfix) > 0 {
				branchingModel.Types[i].Prefix = command.flags.prefixHotfix
			}
		case "RELEASE":
			branchingModel.Types[i].Enabled = command.flags.enableRelease
			if len(command.flags.prefixRelease) > 0 {
				branchingModel.Types[i].Prefix = command.flags.prefixRelease
			}
		default:
			return fmt.Errorf("unsupported branching model type %s", t.Id)
		}
	}

	err = client.SetBranchingModel(
		command.flags.project,
		command.flags.repository,
		branchingModel,
	)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] set branching model for repository %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}
