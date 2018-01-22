// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"fmt"
	"strconv"

	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// ShowPermissionsCommand define base struct for ShowPermissions actions
type ShowPermissionsCommand struct {
	Settings *settings.BitAdminSettings
	flags    *ShowPermissionsFlags
}

// ShowPermissionsFlags define flags required by the ShowPermissionsAction
type ShowPermissionsFlags struct {
	project    string
	repository string
}

// GetCommand provide a ready to use cli.Command
func (command *ShowPermissionsCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "show-permission",
		Usage:  "Show permissions on given repository",
		Action: command.ShowPermissionsAction,
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
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

type OutputRow struct {
	Name       string
	Type       string
	Project    string
	Repository string
	Read       bool
	Write      bool
	Merge      bool
}

type OutputRows []OutputRow

// ShowPermissionsAction display the current user / group permissions on given repository
func (command *ShowPermissionsCommand) ShowPermissionsAction(context *cli.Context) error {

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	userResponse, err := client.GetRepositoryUserPermission(
		command.flags.project,
		command.flags.repository,
		bitclient.GetRepositoryUserPermissionRequest{},
	)

	if err != nil {
		return err
	}

	groupResponse, err := client.GetRepositoryGroupPermission(
		command.flags.project,
		command.flags.repository,
		bitclient.GetRepositoryGroupPermissionRequest{},
	)

	if err != nil {
		return err
	}

	branchRestrictions, err := client.GetRepositoryBranchRestrictions(
		command.flags.project,
		command.flags.repository,
		bitclient.GetRepositoryBranchRestrictionRequest{
			MatcherType: "BRANCH",
			MatcherId:   "refs/heads/master",
			Type:        "read-only",
		},
	)

	if err != nil {
		return err
	}

	var masterRestriction bitclient.BranchRestriction
	if len(branchRestrictions) > 0 {
		masterRestriction = branchRestrictions[0]
	}

	projectUserResponse, err := client.GetProjectUserPermission(command.flags.project, bitclient.GetProjectUserPermissionRequest{})
	if err != nil {
		return err
	}
	projectGroupResponse, err := client.GetProjectGroupPermission(command.flags.project, bitclient.GetProjectGroupPermissionRequest{})
	if err != nil {
		return err
	}

	var rowList OutputRows

	for _, userPermission := range userResponse.Values {
		rowList = rowList.appendIfNew(OutputRow{
			Name:       userPermission.User.Slug,
			Type:       "user",
			Project:    command.flags.project,
			Repository: command.flags.repository,
			Read:       hasRead(userPermission.Permission),
			Write:      hasWrite(userPermission.Permission),
			Merge:      hasMerge(userPermission.User.Slug, masterRestriction),
		})
	}

	for _, groupPermission := range groupResponse.Values {
		rowList = rowList.appendIfNew(OutputRow{
			Name:       groupPermission.Group.Name,
			Type:       "group",
			Project:    command.flags.project,
			Repository: command.flags.repository,
			Read:       hasRead(groupPermission.Permission),
			Write:      hasWrite(groupPermission.Permission),
			Merge:      hasMerge(groupPermission.Group.Name, masterRestriction),
		})
	}

	for _, mergeUser := range masterRestriction.Users {
		rowList = rowList.appendIfNew(OutputRow{
			Name:       mergeUser.Slug,
			Type:       "user",
			Project:    command.flags.project,
			Repository: command.flags.repository,
			Read:       true,
			Write:      true,
			Merge:      true,
		})
	}

	for _, mergeGroup := range masterRestriction.Groups {
		rowList = rowList.appendIfNew(OutputRow{
			Name:       mergeGroup,
			Type:       "group",
			Project:    command.flags.project,
			Repository: command.flags.repository,
			Read:       true,
			Write:      true,
			Merge:      true,
		})
	}

	for _, perm := range projectUserResponse.Values {
		rowList = rowList.appendIfNew(OutputRow{
			Name:       perm.User.Slug,
			Type:       "user",
			Project:    command.flags.project,
			Repository: command.flags.repository,
			Read:       hasRead(perm.Permission),
			Write:      hasWrite(perm.Permission),
			Merge:      false,
		})
	}

	for _, perm := range projectGroupResponse.Values {
		rowList = rowList.appendIfNew(OutputRow{
			Name:       perm.Group.Name,
			Type:       "group",
			Project:    command.flags.project,
			Repository: command.flags.repository,
			Read:       hasRead(perm.Permission),
			Write:      hasWrite(perm.Permission),
			Merge:      false,
		})
	}

	fmt.Printf("%s\n", rowList)

	return nil
}

func (rows OutputRows) appendIfNew(row OutputRow) OutputRows {
	for _, r := range rows {
		if r.Name == row.Name {
			rows = append(rows, row)
		}
	}

	return rows
}

func (rows OutputRows) String() string {
	out := "Name;Type;Project;Repository;Read;Write;Merge"
	for _, row := range rows {
		out += fmt.Sprintf("%s", row)
	}

	return out
}

func (row OutputRow) String() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s\n",
		row.Name,
		row.Type,
		row.Project,
		row.Repository,
		strconv.FormatBool(row.Read),
		strconv.FormatBool(row.Write),
		strconv.FormatBool(row.Merge),
	)
}

func hasRead(permission string) bool {
	var hasRead bool

	if write := hasWrite(permission); write == true {
		hasRead = true
	} else {
		switch permission {
		case bitclient.REPO_READ:
		case bitclient.PROJECT_READ:
			hasRead = true
		default:
			hasRead = false
		}
	}
	return hasRead
}

func hasWrite(permission string) bool {
	var hasWrite bool

	switch permission {
	case bitclient.REPO_ADMIN:
	case bitclient.REPO_WRITE:
	case bitclient.PROJECT_WRITE:
		hasWrite = true
	default:
		hasWrite = false
	}

	return hasWrite
}

func hasMerge(slug string, branchRestriction bitclient.BranchRestriction) bool {

	for _, u := range branchRestriction.Users {
		if slug == u.Slug {
			return true
		}
	}

	for _, g := range branchRestriction.Groups {
		if slug == g {
			return true
		}
	}

	return false
}
