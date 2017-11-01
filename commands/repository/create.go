// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
)

// CreateCommand define base struct for Create actions
type CreateCommand struct {
	Settings *settings.BitAdminSettings
	flags    *CreateCommandFlags
}

// CreateCommandFlags hold flag values for the CreateCommand
type CreateCommandFlags struct {
	project  string
	name     string
	scm      string
	forkable bool
}

// GetCommand provide a ready to use cli.Command
func (command *CreateCommand) GetCommand(fileCache *helper.FileCache) cli.Command {
	return cli.Command{
		Name:   "create",
		Usage:  "Create a new repository",
		Action: command.CreateRepositoryAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Usage:       "The `<project_key>` where the repository will be created",
				Destination: &command.flags.project,
			},
			cli.StringFlag{
				Name:        "name",
				Usage:       "The `<repository_name>` to create",
				Destination: &command.flags.name,
			},
			cli.StringFlag{
				Name:        "scm",
				Value:       "git",
				Usage:       "The `<scm>` to use.",
				Destination: &command.flags.scm,
			},
			cli.BoolFlag{
				Name:        "forkable",
				Usage:       "Allow a repository to be forked",
				Destination: &command.flags.forkable,
			},
		},
		BashComplete: func(c *cli.Context) {
			helper.AutoComplete(c, fileCache)
		},
	}
}

// CreateRepositoryAction use flag values to create a new repository
func (command *CreateCommand) CreateRepositoryAction(context *cli.Context) error {
	if len(command.flags.project) == 0 {
		return fmt.Errorf("flag --project is required")
	}

	if len(command.flags.name) == 0 {
		return fmt.Errorf("flag --name is required")
	}

	if len(command.flags.scm) == 0 {
		return fmt.Errorf("flag --scm is required")
	}

	requestData := bitclient.CreateRepositoryRequest{
		Name:     command.flags.name,
		ScmId:    command.flags.scm,
		Forkable: command.flags.forkable,
	}

	client, err := command.Settings.GetAPIClient()

	if err != nil {
		return err
	}

	resp, err := client.CreateRepository(command.flags.project, requestData)

	if err != nil {
		switch terr := err.(type) {
		case bitclient.RequestError:
			switch terr.Code {
			case 404:
				return fmt.Errorf("project {%s} does not exists", command.flags.project)
			case 409:
				return fmt.Errorf("repository {%s} already exists", command.flags.name)
			}
		}
		return err
	}

	fmt.Println("Repository created")
	fmt.Println("Quick links :")
	helper.PrintLinks(resp.Links)
	fmt.Println()

	fileCache := command.Settings.GetFileCache()
	fileCache.Repositories = append(fileCache.Repositories, resp)
	fileCache.Save()
	return nil
}
