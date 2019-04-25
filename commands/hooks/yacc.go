// Package hooks hold actions on the Bitbucket hooks
package hooks

import (
	"errors"
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
	"encoding/json"
	"github.com/daeMOn63/bitclient"
	"github.com/google/go-cmp/cmp"
	"strings"
	"unicode"
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

	yaccHookSettingsCommand := YaccHookSettingsCommand{
		Settings: command.Settings,
		flags:    &YaccHookGetSettingsCommandFlags{},
	}

	yaccDiffHookSettingsCommand := YaccHookDiffSettingsCommand{
		Settings: command.Settings,
		flags:    &YaccHookDiffSettingsCommandFlags{},
	}

	return cli.Command{
		Name:  "yet-another-commit-checker",
		Usage: "Yet Another Commit Checker hook operations",
		Subcommands: []cli.Command{
			yaccEnableCommand.GetCommand(),
			yaccDisableCommand.GetCommand(),
			yaccHookSettingsCommand.GetCommand(),
			yaccDiffHookSettingsCommand.GetCommand(),
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
				Usage:       "If present, JIRA issues must match this JQL query.",
				Destination: &command.flags.settings.IssueJqlMatcher,
			},
			cli.StringFlag{
				Name:        "branchNameRegex",
				Usage:       "If present, only branches with names that match this regex will be allowed to be created.",
				Destination: &command.flags.settings.BranchNameRegex,
			},
			cli.StringFlag{
				Name:        "errorMessageHeader",
				Usage:       "If present, the default error message header will be replaced by this text.",
				Destination: &command.flags.settings.ErrorMessageHeader,
			},
			cli.StringFlag{
				Name:        "errorMessageCommiterEmail",
				Usage:       "If present, this text will be shown when the Require Matching Committer Email check fails.",
				Destination: &command.flags.settings.ErrorMessageCommiterEmail,
			},
			cli.StringFlag{
				Name:        "errorMessageCommiterEmailRegex",
				Usage:       "If present, this text will be shown when the Committer Email Regex check fails.",
				Destination: &command.flags.settings.ErrorMessageCommiterEmailRegex,
			},
			cli.StringFlag{
				Name:        "errorMessageCommiterName",
				Usage:       "If present, this text will be shown when the Require Matching Committer Name check fails.",
				Destination: &command.flags.settings.ErrorMessageCommiterName,
			},
			cli.StringFlag{
				Name:        "errorMessageCommitRegex",
				Usage:       "If present, this text will be shown when the Commit Message Regex check fails.",
				Destination: &command.flags.settings.ErrorMessageCommitRegex,
			},
			cli.StringFlag{
				Name:        "errorMessageIssueJQL",
				Usage:       "If present, this text will be shown when the Issue Jql Matcher check fails.",
				Destination: &command.flags.settings.ErrorMessageIssueJQL,
			},
			cli.StringFlag{
				Name:        "errorMessageBranchName",
				Usage:       "If present, this text will be shown when the Branch Name Regex check fails.",
				Destination: &command.flags.settings.ErrorMessageBranchName,
			},
			cli.StringFlag{
				Name:        "errorMessageFooter",
				Usage:       "If present, this text will be included at the end of the YACC error message.",
				Destination: &command.flags.settings.ErrorMessageFooter,
			},
			cli.BoolFlag{
				Name:        "excludeMergeCommits",
				Usage:       "Exclude merge commits from commit requirements.",
				Destination: &command.flags.settings.ExcludeMergeCommits,
			},
			cli.StringFlag{
				Name:        "excludeByRegex",
				Usage:       "Exclude commits if commit message matches this regex.",
				Destination: &command.flags.settings.ExcludeByRegex,
			},
			cli.StringFlag{
				Name:        "excludeBranchRegex",
				Usage:       "Exclude commits to branches matching this regex.",
				Destination: &command.flags.settings.ExcludeBranchRegex,
			},
			cli.BoolFlag{
				Name:        "excludeServiceUserCommits",
				Usage:       "Exclude commits from service users with access keys (e.g. CI Server) from commit requirements.",
				Destination: &command.flags.settings.ExcludeServiceUserCommits,
			},
			cli.StringFlag{
				Name:        "excludeUsers",
				Usage:       "Exclude commits from users. Separate multiple user names with a comma.",
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

// YaccHookSettingsCommand define the command to get setting for YACC hook
type YaccHookSettingsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *YaccHookGetSettingsCommandFlags
}

// YaccHookGetSettingsCommandFlags define the flags of the YaccHookSettingsCommand
type YaccHookGetSettingsCommandFlags struct {
	project         string
	repository      string
	defaultSettings bitclient.YaccHookSettings
}

// GetCommand provide a ready to use cli.Command
func (command *YaccHookSettingsCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "get-settings",
		Usage:  "Get settings for Yet Another Commit Checker hook",
		Action: command.GetSettings,
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

// GetSettings contains logic to get current yacc settings
func (command *YaccHookSettingsCommand) GetSettings(context *cli.Context) error {
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

	yacc, err := client.GetYACCHookSettings(
		command.flags.project,
		command.flags.repository,
	)

	data, err := json.Marshal(yacc)
	if err != nil {
		return err
	}

	fmt.Printf("%s", data)

	return nil
}


// YaccHookSettingsCommand define the command to get setting for YACC hook
type YaccHookDiffSettingsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *YaccHookDiffSettingsCommandFlags
}

// YaccHookGetSettingsCommandFlags define the flags of the YaccHookSettingsCommand
type YaccHookDiffSettingsCommandFlags struct {
	project         string
	repository      string
	defaultSettings bitclient.YaccHookSettings
}

// GetCommand provide a ready to use cli.Command
func (command *YaccHookDiffSettingsCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "diff-settings",
		Usage:  "Diff YACC settings of a repository against default settings",
		Action: command.DiffSettings,
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
				Name:        "default-requireMatchingAuthorEmail",
				Usage:       "Require that the commit committer's email matches the Stash user's email.",
				Destination: &command.flags.defaultSettings.RequireMatchingAuthorEmail,
			},
			cli.BoolFlag{
				Name:        "default-requireMatchingAuthorName",
				Usage:       "Require that the commit committer's name matches the Stash user's name.",
				Destination: &command.flags.defaultSettings.RequireMatchingAuthorName,
			},
			cli.StringFlag{
				Name:        "default-committerEmailRegex",
				Usage:       "Require that commit email match this regular expression.",
				Destination: &command.flags.defaultSettings.CommitterEmailRegex,
			},
			cli.StringFlag{
				Name:        "default-commitMessageRegex",
				Usage:       "Require that commit messages match this regular expression.",
				Destination: &command.flags.defaultSettings.CommitMessageRegex,
			},
			cli.BoolFlag{
				Name:        "default-requireJiraIssue",
				Usage:       "Require that the commit message contains valid JIRA issue(s).",
				Destination: &command.flags.defaultSettings.RequireJiraIssue,
			},
			cli.BoolFlag{
				Name:        "default-ignoreUnknownIssueProjectKeys",
				Usage:       "Items in the commit message that do not contain a valid JIRA project key (such as UTF-8) will be ignored.",
				Destination: &command.flags.defaultSettings.IgnoreUnknownIssueProjectKeys,
			},
			cli.StringFlag{
				Name:        "default-issueJqlMatcher",
				Usage:       "If present, JIRA issues must match this JQL query.",
				Destination: &command.flags.defaultSettings.IssueJqlMatcher,
			},
			cli.StringFlag{
				Name:        "default-branchNameRegex",
				Usage:       "If present, only branches with names that match this regex will be allowed to be created.",
				Destination: &command.flags.defaultSettings.BranchNameRegex,
			},
			cli.StringFlag{
				Name:        "default-errorMessageHeader",
				Usage:       "If present, the default error message header will be replaced by this text.",
				Destination: &command.flags.defaultSettings.ErrorMessageHeader,
			},
			cli.StringFlag{
				Name:        "default-errorMessageCommiterEmail",
				Usage:       "If present, this text will be shown when the Require Matching Committer Email check fails.",
				Destination: &command.flags.defaultSettings.ErrorMessageCommiterEmail,
			},
			cli.StringFlag{
				Name:        "default-errorMessageCommiterEmailRegex",
				Usage:       "If present, this text will be shown when the Committer Email Regex check fails.",
				Destination: &command.flags.defaultSettings.ErrorMessageCommiterEmailRegex,
			},
			cli.StringFlag{
				Name:        "default-errorMessageCommiterName",
				Usage:       "If present, this text will be shown when the Require Matching Committer Name check fails.",
				Destination: &command.flags.defaultSettings.ErrorMessageCommiterName,
			},
			cli.StringFlag{
				Name:        "default-errorMessageCommitRegex",
				Usage:       "If present, this text will be shown when the Commit Message Regex check fails.",
				Destination: &command.flags.defaultSettings.ErrorMessageCommitRegex,
			},
			cli.StringFlag{
				Name:        "default-errorMessageIssueJQL",
				Usage:       "If present, this text will be shown when the Issue Jql Matcher check fails.",
				Destination: &command.flags.defaultSettings.ErrorMessageIssueJQL,
			},
			cli.StringFlag{
				Name:        "default-errorMessageBranchName",
				Usage:       "If present, this text will be shown when the Branch Name Regex check fails.",
				Destination: &command.flags.defaultSettings.ErrorMessageBranchName,
			},
			cli.StringFlag{
				Name:        "default-errorMessageFooter",
				Usage:       "If present, this text will be included at the end of the YACC error message.",
				Destination: &command.flags.defaultSettings.ErrorMessageFooter,
			},
			cli.BoolFlag{
				Name:        "default-excludeMergeCommits",
				Usage:       "Exclude merge commits from commit requirements.",
				Destination: &command.flags.defaultSettings.ExcludeMergeCommits,
			},
			cli.StringFlag{
				Name:        "default-excludeByRegex",
				Usage:       "Exclude commits if commit message matches this regex.",
				Destination: &command.flags.defaultSettings.ExcludeByRegex,
			},
			cli.StringFlag{
				Name:        "default-excludeBranchRegex",
				Usage:       "Exclude commits to branches matching this regex.",
				Destination: &command.flags.defaultSettings.ExcludeBranchRegex,
			},
			cli.BoolFlag{
				Name:        "default-excludeServiceUserCommits",
				Usage:       "Exclude commits from service users with access keys (e.g. CI Server) from commit requirements.",
				Destination: &command.flags.defaultSettings.ExcludeServiceUserCommits,
			},
			cli.StringFlag{
				Name:        "default-excludeUsers",
				Usage:       "Exclude commits from users. Separate multiple user names with a comma.",
				Destination: &command.flags.defaultSettings.ExcludeUsers,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}


// DiffSettings contains logic to get diff between defined default settings and current setting got for specific repository
func (command *YaccHookDiffSettingsCommand) DiffSettings(context *cli.Context) error {
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

	yacc, err := client.GetYACCHookSettings(
		command.flags.project,
		command.flags.repository,
	)

	data, err := json.Marshal(yacc)
	if err != nil {
		return err
	}

	if command.flags.defaultSettings == (bitclient.YaccHookSettings{}) {
		fmt.Println("Default settings empty, no comparison will be made")

		return nil
	}

	isCorrect  := "NO"
	isOK, diff := isYACCSettingsCorrectWithDiff(command.flags.defaultSettings,yacc)

	if isOK == true {
		isCorrect = "YES"
	}

	fmt.Printf("%s# %s# %s", isCorrect, data, diff.String())

	return nil
}

// isYACCSettingsCorrectWithDiff Compares YACC hook settings for a repository and check whether these are matching the desired default settings
func isYACCSettingsCorrectWithDiff(defaultYACCSettings, yaccSettings bitclient.YaccHookSettings) (bool, DiffReporter) {

	opts := cmp.Transformer("StripWhitespace", func(x bitclient.YaccHookSettings) bitclient.YaccHookSettings {
		temp := x
		temp.IssueJqlMatcher = SpaceStringsBuilder(x.IssueJqlMatcher)
		temp.ErrorMessageIssueJQL = SpaceStringsBuilder(x.ErrorMessageIssueJQL)

		return temp
	})

	var r DiffReporter
	cmp.Diff(defaultYACCSettings, yaccSettings, cmp.Reporter(&r))

	return cmp.Equal(defaultYACCSettings, yaccSettings, opts), r
}

// SpaceStringsBuilder Strips whitespace from a string
func SpaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// DiffReporter is a simple custom reporter that only records differences
// detected during comparison.
type DiffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		r.diffs = append(r.diffs, fmt.Sprintf("%v:-: %v+: %v\n", r.path, vx, vy))
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}