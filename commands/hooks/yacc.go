// Package hooks hold actions on the Bitbucket hooks
package hooks

import (
	"errors"
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

const yaccHookKey = "com.isroot.stash.plugin.yacc:yaccHook"

// YaccHookCommand define the command for the Yet Another Commit Checker hook
type YaccHookCommand struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *YaccHookCommand) GetCommand() cli.Command {

	yaccEnableCommand := YaccHookEnableCommand{
		Settings: command.Settings,
		flags:    &YaccHookEnableCommandFlags{},
	}

	yaccDisableCommand := YaccHookDisableCommand{
		Settings: command.Settings,
		flags:    &YaccHookDisableCommandFlags{},
	}

	return cli.Command{
		Name:  "yet-another-commit-checker",
		Usage: "Yet Another Commit Checker hook operations",
		Subcommands: []cli.Command{
			yaccEnableCommand.GetCommand(),
			yaccDisableCommand.GetCommand(),
		},
	}
}

// YaccHookEnableCommand define the command to enable the YACC hook
type YaccHookEnableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *YaccHookEnableCommandFlags
}

// YaccSettings define the settings of the YACC hook
type YaccSettings struct {
	RequireMatchingAuthorEmail     bool   `json:"requireMatchingAuthorEmail,omitempty"`
	RequireMatchingAuthorName      bool   `json:"requireMatchingAuthorName,omitempty"`
	CommitterEmailRegex            string `json:"committerEmailRegex,omitempty"`
	CommitMessageRegex             string `json:"commitMessageRegex,omitempty"`
	RequireJiraIssue               bool   `json:"requireJiraIssue,omitempty"`
	IgnoreUnknownIssueProjectKeys  bool   `json:"ignoreUnknownIssueProjectKeys,omitempty"`
	IssueJqlMatcher                string `json:"issueJqlMatcher,omitempty"`
	BranchNameRegex                string `json:"branchNameRegex,omitempty"`
	ErrorMessageHeader             string `json:"errorMessageHeader,omitempty"`
	ErrorMessageCommiterEmail      string `json:"errorMessage.COMMITER_EMAIL,omitempty"`
	ErrorMessageCommiterEmailRegex string `json:"errorMessage.COMMITER_EMAIL_REGEX,omitempty"`
	ErrorMessageCommiterName       string `json:"errorMessage.COMMITER_NAME,omitempty"`
	ErrorMessageCommitRegex        string `json:"errorMessage.COMMIT_REGEX,omitempty"`
	ErrorMessageIssueJQL           string `json:"errorMessage.ISSUE_JQL,omitempty"`
	ErrorMessageBranchName         string `json:"errorMessage.BRANCH_NAME,omitempty"`
	ErrorMessageFooter             string `json:"errorMessageFooter,omitempty"`
	ExcludeMergeCommits            bool   `json:"excludeMergeCommits,omitempty"`
	ExcludeByRegex                 string `json:"excludeByRegex,omitempty"`
	ExcludeBranchRegex             string `json:"excludeBranchRegex,omitempty"`
	ExcludeServiceUserCommits      bool   `json:"excludeServiceUserCommits,omitempty"`
	ExcludeUsers                   string `json:"excludeUsers,omitempty"`
}

// YaccHookEnableCommandFlags define the flags for the YaccHookEnableCommand
type YaccHookEnableCommandFlags struct {
	project    string
	repository string
	settings   YaccSettings
}

// GetCommand provide a ready to use cli.Command
func (command *YaccHookEnableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "enable",
		Usage:  "Enable Yet Another Commit Checker hook and set its configuration",
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
			cli.BoolFlag{
				Name:        "requireMatchingAuthorEmail",
				Usage:       "Require that the commit committer's email matches the Stash user's email.",
				Destination: &command.flags.settings.RequireMatchingAuthorEmail,
			},
			cli.BoolFlag{
				Name:        "requireMatchingAuthorName",
				Usage:       "Require that the commit committer's name matches the Stash user's name.",
				Destination: &command.flags.settings.RequireMatchingAuthorName,
			},
			cli.StringFlag{
				Name:        "committerEmailRegex",
				Usage:       "Require that commit email match this regular expression.",
				Destination: &command.flags.settings.CommitterEmailRegex,
			},
			cli.StringFlag{
				Name:        "commitMessageRegex",
				Usage:       "Require that commit messages match this regular expression.",
				Destination: &command.flags.settings.CommitMessageRegex,
			},
			cli.BoolFlag{
				Name:        "requireJiraIssue",
				Usage:       "Require that the commit message contains valid JIRA issue(s).",
				Destination: &command.flags.settings.RequireJiraIssue,
			},
			cli.BoolFlag{
				Name:        "ignoreUnknownIssueProjectKeys",
				Usage:       "Items in the commit message that do not contain a valid JIRA project key (such as UTF-8) will be ignored.",
				Destination: &command.flags.settings.IgnoreUnknownIssueProjectKeys,
			},
			cli.StringFlag{
				Name:        "issueJqlMatcher",
				Usage:       "TODO",
				Destination: &command.flags.settings.IssueJqlMatcher,
			},
			cli.StringFlag{
				Name:        "branchNameRegex",
				Usage:       "TODO",
				Destination: &command.flags.settings.BranchNameRegex,
			},
			cli.StringFlag{
				Name:        "errorMessageHeader",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageHeader,
			},
			cli.StringFlag{
				Name:        "errorMessageCommiterEmail",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageCommiterEmail,
			},
			cli.StringFlag{
				Name:        "errorMessageCommiterEmailRegex",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageCommiterEmailRegex,
			},
			cli.StringFlag{
				Name:        "errorMessageCommiterName",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageCommiterName,
			},
			cli.StringFlag{
				Name:        "errorMessageCommitRegex",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageCommitRegex,
			},
			cli.StringFlag{
				Name:        "errorMessageIssueJQL",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageIssueJQL,
			},
			cli.StringFlag{
				Name:        "errorMessageBranchName",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageBranchName,
			},
			cli.StringFlag{
				Name:        "errorMessageFooter",
				Usage:       "TODO",
				Destination: &command.flags.settings.ErrorMessageFooter,
			},
			cli.BoolFlag{
				Name:        "excludeMergeCommits",
				Usage:       "TODO",
				Destination: &command.flags.settings.ExcludeMergeCommits,
			},
			cli.StringFlag{
				Name:        "excludeByRegex",
				Usage:       "TODO",
				Destination: &command.flags.settings.ExcludeByRegex,
			},
			cli.StringFlag{
				Name:        "excludeBranchRegex",
				Usage:       "TODO",
				Destination: &command.flags.settings.ExcludeBranchRegex,
			},
			cli.BoolFlag{
				Name:        "excludeServiceUserCommits",
				Usage:       "TODO",
				Destination: &command.flags.settings.ExcludeServiceUserCommits,
			},
			cli.StringFlag{
				Name:        "excludeUsers",
				Usage:       "TODO",
				Destination: &command.flags.settings.ExcludeUsers,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// EnableAction contains logic to turn on the hook and set its configuration
func (command *YaccHookEnableCommand) EnableAction(context *cli.Context) error {
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
		yaccHookKey,
		command.flags.settings,
	)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Enabled and configured yacc hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}

// YaccHookDisableCommand define the command to disable the YACC hook
type YaccHookDisableCommand struct {
	Settings *settings.BitAdminSettings
	flags    *YaccHookDisableCommandFlags
}

// YaccHookDisableCommandFlags define the flags of the YaccHookDisableCommand
type YaccHookDisableCommandFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *YaccHookDisableCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "disable",
		Usage:  "Disable Yet Another Commit Checker hook",
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
func (command *YaccHookDisableCommand) DisableAction(context *cli.Context) error {
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

	err = client.DisableHook(command.flags.project, command.flags.repository, yaccHookKey)

	if err != nil {
		return err
	}

	fmt.Printf("[OK] Disabled yacc hook on %s/%s\n", command.flags.project, command.flags.repository)

	return nil
}
