// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"errors"
	"fmt"

	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// SetDefaultReviewersCommand define base struct for SetDefaultReviewer actions
type SetDefaultReviewersCommand struct {
	Settings *settings.BitAdminSettings
	flags    *SetDefaultReviewersCommandFlags
}

// SetDefaultReviewersCommandFlags hold flag values for the SetDefaultReviewerCommand
type SetDefaultReviewersCommandFlags struct {
	project           string
	repository        string
	usernames         cli.StringSlice
	branchRef         string
	requiredApprovers uint
	replace           bool
}

// GetCommand provide a ready to use cli.Command
func (command *SetDefaultReviewersCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "set-default-reviewers",
		Usage:  "Set default reviewers on given repository",
		Action: command.SetDefaultReviewersAction,
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
			cli.StringSliceFlag{
				Name:  "username",
				Usage: "The `<username>` to be added on the repository. Can be repeated multiple times",
				Value: &command.flags.usernames,
			},
			cli.UintFlag{
				Name:        "requiredApprovers",
				Usage:       "`<requiredApprovers>` set the minimum number of approval required by default reviewers",
				Destination: &command.flags.requiredApprovers,
			},
			cli.StringFlag{
				Name:        "branchRef",
				Usage:       "`<branchRef>` define the pull request's target branch where default reviewers will be set (ie: refs/heads/master)",
				Destination: &command.flags.branchRef,
			},
			cli.BoolFlag{
				Name:        "replace",
				Usage:       "Setting this flag will replace existing default reviewers by provided ones.",
				Destination: &command.flags.replace,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, command.Settings.GetFileCache())
		},
	}
}

// SetDefaultReviewersAction allow to set the default reviewers on given repository.
func (command *SetDefaultReviewersCommand) SetDefaultReviewersAction(context *cli.Context) error {
	client, err := command.Settings.GetAPIClient()
	if err != nil {
		return err
	}

	cache := command.Settings.GetFileCache()

	if len(command.flags.repository) <= 0 {
		return errors.New("--repository flag is required")
	}
	if len(command.flags.project) <= 0 {
		return errors.New("--project flag is required")
	}
	if len(command.flags.branchRef) <= 0 {
		return errors.New("--branchRef flag is required")
	}
	if len(command.flags.usernames) == 0 {
		return fmt.Errorf("At least one --username is required")
	}

	repo, err := cache.FindRepository(command.flags.project, command.flags.repository)
	if err != nil {
		return err
	}

	var users []bitclient.User

	for _, username := range command.flags.usernames {

		user, err := cache.FindUserByUsername(username)
		if err != nil {
			return err
		}

		users = append(users, user)
	}

	settings, err := client.GetRepositoryDefaultReviewers(command.flags.project, command.flags.repository)

	exists := false

	for _, setting := range settings {
		if setting.ToRefMatcher.Id == command.flags.branchRef {

			if command.flags.replace == false {

				for _, user := range users {
					found := false
					for _, revUser := range setting.Reviewers {
						if revUser.Slug == user.Slug {
							found = true
							break
						}
					}

					if found == false {
						setting.Reviewers = append(setting.Reviewers, user)
					}
				}
			} else {
				setting.Reviewers = users
			}

			client.UpdateRepositoryDefaultReviewers(command.flags.project, command.flags.repository, setting)
			exists = true
			break
		}
	}

	// No default reviewers exists for given branch, create it
	if exists == false {
		setting := bitclient.DefaultReviewers{
			Repository: repo,
			FromRefMatcher: bitclient.Matcher{
				Id:   "ANY_REF_MATCHER_ID",
				Type: bitclient.MatcherType{Id: "ANY_REF"},
			},
			ToRefMatcher: bitclient.Matcher{
				Id:   command.flags.branchRef,
				Type: bitclient.MatcherType{Id: "BRANCH"},
			},
			RequiredApprovals: int(command.flags.requiredApprovers),
			Reviewers:         users,
		}

		client.CreateRepositoryDefaultReviewers(command.flags.project, command.flags.repository, setting)
	}

	fmt.Printf(
		"Added %d users as default reviewers on %s for %s/%s\n",
		len(users),
		command.flags.branchRef,
		command.flags.project,
		command.flags.repository,
	)

	return nil
}
