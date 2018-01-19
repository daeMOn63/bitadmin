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
		//fmt.Printf("user %s - %s\n", userPermission.User.Slug, userPermission.Permission)
		rowList = append(rowList, OutputRow{
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
		//fmt.Printf("group %s - %s\n", groupPermission.Group.Name, groupPermission.Permission)
		rowList = append(rowList, OutputRow{
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
		if rowList.containsName(mergeUser.Slug) == false {
			rowList = append(rowList, OutputRow{
				Name:       mergeUser.Slug,
				Type:       "user",
				Project:    command.flags.project,
				Repository: command.flags.repository,
				Read:       true,
				Write:      true,
				Merge:      true,
			})
		}
	}
	for _, mergeGroup := range masterRestriction.Groups {
		if rowList.containsName(mergeGroup) == false {
			rowList = append(rowList, OutputRow{
				Name:       mergeGroup,
				Type:       "group",
				Project:    command.flags.project,
				Repository: command.flags.repository,
				Read:       true,
				Write:      true,
				Merge:      true,
			})
		}
	}

	for _, perm := range projectUserResponse.Values {
		if rowList.containsName(perm.User.Slug) == false {
			rowList = append(rowList, OutputRow{
				Name:       perm.User.Slug,
				Type:       "user",
				Project:    command.flags.project,
				Repository: command.flags.repository,
				Read:       hasRead(perm.Permission),
				Write:      hasWrite(perm.Permission),
				Merge:      false,
			})
		}
	}

	for _, perm := range projectGroupResponse.Values {
		if rowList.containsName(perm.Group.Name) == false {
			rowList = append(rowList, OutputRow{
				Name:       perm.Group.Name,
				Type:       "group",
				Project:    command.flags.project,
				Repository: command.flags.repository,
				Read:       hasRead(perm.Permission),
				Write:      hasWrite(perm.Permission),
				Merge:      false,
			})
		}
	}

	fmt.Printf("%s\n", rowList)

	return nil
}

func (rows OutputRows) containsName(name string) bool {
	for _, r := range rows {
		if r.Name == name {
			return true
		}
	}

	return false
}

func (rows OutputRows) String() string {
	var out string
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
	if permission == bitclient.REPO_ADMIN || permission == bitclient.REPO_WRITE || permission == bitclient.REPO_READ || permission == bitclient.PROJECT_READ || permission == bitclient.PROJECT_WRITE {
		return true
	}

	return false
}

func hasWrite(permission string) bool {
	if permission == bitclient.REPO_ADMIN || permission == bitclient.REPO_WRITE || permission == bitclient.PROJECT_WRITE {
		return true
	}

	return false
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
