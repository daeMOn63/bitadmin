// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"errors"
	"fmt"

	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// SonarCommand define base struct for SonarCommand actions
type SonarCommand struct {
	Settings *settings.BitAdminSettings
	flags    *SonarCommandFlags
}

// SonarCommandFlags define flags required by the SonarAction
type SonarCommandFlags struct {
	project                      string
	repository                   string
	enabled                      bool
	sonarMasterProjectKey        string
	sonarProjectBaseKey          string
	analysisMode                 string
	useSonarBranchFeature        bool
	showIssuesInSource           bool
	showOnlyNewOrChangedLines    bool
	illegalBranchCharReplacement string
	projectCleanupEnabled        bool
	// TODO: add more when needed :)
}

// GetCommand provide a ready to use cli.Command
func (command *SonarCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "sonar",
		Usage:  "Update sonar setting for a given repository",
		Action: command.SonarAction,
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
			cli.BoolFlag{
				Name:        "enabled",
				Usage:       "Enable the sonar plugin",
				Destination: &command.flags.enabled,
			},
			cli.StringFlag{
				Name:        "sonarMasterProjectKey",
				Usage:       "The sonar master project key to link this repository with",
				Destination: &command.flags.sonarMasterProjectKey,
			},
			cli.StringFlag{
				Name:        "sonarProjectBaseKey",
				Usage:       "The sonar project base key to link this repository with. Not used when --useSonarBranchFeature is set",
				Destination: &command.flags.sonarProjectBaseKey,
			},
			cli.StringFlag{
				Name:        "analysisMode",
				Usage:       "Sonar analysis mode (one of LEAK_PERIOD or BRANCH_DIFF)",
				Destination: &command.flags.analysisMode,
			},
			cli.BoolFlag{
				Name:        "useSonarBranchFeature",
				Usage:       "Enable the new branching feature introduced with commercial editions of SonarQube 6.7",
				Destination: &command.flags.useSonarBranchFeature,
			},
			cli.BoolFlag{
				Name:        "showIssuesInSource",
				Usage:       "Show issues in sources",
				Destination: &command.flags.showIssuesInSource,
			},
			cli.BoolFlag{
				Name:        "showOnlyNewOrChangedLines",
				Usage:       "Show issues only on changed lines",
				Destination: &command.flags.showOnlyNewOrChangedLines,
			},
			cli.StringFlag{
				Name:        "illegalBranchCharReplacement",
				Usage:       "Replace illegal character in branches (for SonarQube <5.0)",
				Destination: &command.flags.analysisMode,
			},
			cli.BoolFlag{
				Name:        "projectCleanupEnabled",
				Usage:       "Enable project cleanup",
				Destination: &command.flags.projectCleanupEnabled,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// SonarAction allow to turn on the sonar cleanup setting on all available repositories
func (command *SonarCommand) SonarAction(context *cli.Context) error {
	client, _ := command.Settings.GetAPIClient()

	if len(command.flags.repository) <= 0 {
		return errors.New("--repository flag is required")
	}
	if len(command.flags.project) <= 0 {
		return errors.New("--project flag is required")
	}
	sonarSettings, _ := client.GetSonarSettings(command.flags.project, command.flags.repository)

	sonarSettings.Project.SonarEnabled = command.flags.enabled

	if len(command.flags.sonarMasterProjectKey) > 0 {
		sonarSettings.Project.MasterProjectKey = command.flags.sonarMasterProjectKey
	}

	if len(command.flags.sonarProjectBaseKey) > 0 {
		sonarSettings.Project.ProjectBaseKey = command.flags.sonarProjectBaseKey
	}

	if len(command.flags.analysisMode) > 0 {
		sonarSettings.Project.AnalysisMode = command.flags.analysisMode
	}

	sonarSettings.Project.UseSonarBranchFeature = command.flags.useSonarBranchFeature
	sonarSettings.Project.ShowIssuesInSource = command.flags.showIssuesInSource
	sonarSettings.Project.ShowOnlyNewOrChangedLines = command.flags.showOnlyNewOrChangedLines

	if len(command.flags.illegalBranchCharReplacement) > 0 {
		sonarSettings.Project.IllegalBranchCharReplacement = command.flags.illegalBranchCharReplacement
	}

	sonarSettings.Project.ProjectCleanupEnabled = command.flags.projectCleanupEnabled

	err := client.SetSonarSettings(command.flags.project, command.flags.repository, sonarSettings)
	if err != nil {
		return err
	}

	fmt.Printf("[OK] Updated sonar cleanup settings for repository %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}
