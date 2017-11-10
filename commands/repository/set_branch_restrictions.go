// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"errors"
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
	"strings"
)

// SetBranchRestrictionCommand define base struct for SetBranchRestriction actions
type SetBranchRestrictionCommand struct {
	Settings *settings.BitAdminSettings
	flags    *SetBranchRestrictionCommandFlags
}

// SetBranchRestrictionCommandFlags hold flag values for the SetBranchRestrictionCommand
type SetBranchRestrictionCommandFlags struct {
	project     string
	repository  string
	update      bool
	restriction string
	branchRef   string
	usernames   cli.StringSlice
	groups      cli.StringSlice
}

// GetCommand provide a ready to use cli.Command
func (command *SetBranchRestrictionCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "set-branch-restriction",
		Usage:  "Set branch restrictions on given repository",
		Action: command.SetBranchRestrictionAction,
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
				Name:        "update",
				Usage:       "When set, the current settings won't get overwritten.",
				Destination: &command.flags.update,
			},
			cli.StringFlag{
				Name:        "restriction",
				Usage:       "The `<restriction>` type to set, can be one of  'read-only', 'no-deletes', 'fast-forward-only' or 'pull-request-only'",
				Destination: &command.flags.restriction,
			},
			cli.StringFlag{
				Name:        "branchRef",
				Usage:       "The `<branchRef>` to set the restriction on (ie: refs/heads/master)",
				Destination: &command.flags.branchRef,
			},
			cli.StringSliceFlag{
				Name:  "username",
				Usage: "The `<username>` to be added on the restriction. Can be repeated multiple times",
				Value: &command.flags.usernames,
			},
			cli.StringSliceFlag{
				Name:  "group",
				Usage: "The `<group>` to be added on the restriction. Can be repeated multiple times",
				Value: &command.flags.groups,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// SetBranchRestrictionAction use flag values tp set the branch restrictions on given repository
func (command *SetBranchRestrictionCommand) SetBranchRestrictionAction(context *cli.Context) error {

	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	if len(command.flags.repository) <= 0 {
		return errors.New("--repository flag is required")
	}
	if len(command.flags.project) <= 0 {
		return errors.New("--project flag is required")
	}
	if len(command.flags.restriction) <= 0 {
		return errors.New("--restriction flag is required")
	}
	if len(command.flags.branchRef) <= 0 {
		return errors.New("--branchRef flag is required")
	}

	splittedBranch := strings.Split(command.flags.branchRef, "/")
	displayID := splittedBranch[len(splittedBranch)-1]

	newRestriction := bitclient.SetRepositoryBranchRestrictionsRequest{
		Type: command.flags.restriction,
		Matcher: bitclient.Matcher{
			Id:        command.flags.branchRef,
			DisplayId: displayID,
			Active:    true,
			Type:      bitclient.MatcherType{Id: "BRANCH", Name: "Branch"},
		},
		Users:  command.flags.usernames,
		Groups: command.flags.groups,
	}

	// Set newRestriction to current restriction values if it exists and update is requested.
	if command.flags.update == true {

		getRequestParams := bitclient.GetRepositoryBranchRestrictionRequest{Type: command.flags.restriction}
		restrictions, err := client.GetRepositoryBranchRestrictions(command.flags.project, command.flags.repository, getRequestParams)
		if err != nil {
			return err
		}

		if len(restrictions) == 1 {
			restriction := restrictions[0]

			var origUserSlugs []string
			for _, u := range restriction.Users {
				// We can remove inactive user accounts
				if u.Active == true {
					origUserSlugs = append(origUserSlugs, u.Slug)
				}
			}

			// Keep same matcher
			newRestriction.Id = restriction.Id
			newRestriction.Matcher = restriction.Matcher
			newRestriction.Users = merge(newRestriction.Users, origUserSlugs)
			newRestriction.Groups = merge(newRestriction.Groups, restriction.Groups)
		}

	}

	err = client.SetRepositoryBranchRestrictions(command.flags.project, command.flags.repository, newRestriction)
	if err != nil {
		return err
	}

	action := "updating"
	if command.flags.update == false {
		action = "replacing"
	}

	fmt.Printf(
		"[OK] %s %s restriction on branch %s of %s/%s\n",
		action,
		command.flags.restriction,
		command.flags.branchRef,
		command.flags.project,
		command.flags.repository,
	)

	return nil
}

// merge two string slices removing duplicates values
func merge(s1, s2 []string) []string {

	r := append([]string(nil), s1...)

	for _, s := range s2 {
		exists := false
		for _, e := range r {
			if s == e {
				exists = true
				break
			}
		}

		if exists == false {
			r = append(r, s)
		}
	}

	return r
}
