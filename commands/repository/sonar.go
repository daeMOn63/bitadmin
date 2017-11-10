// Package repository hold actions on the Bitbucket repositories
package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

// EnableSonarCleanupCommand define base struct for EnableSonarCleanup actions
type EnableSonarCleanupCommand struct {
	Settings *settings.BitAdminSettings
}

// GetCommand provide a ready to use cli.Command
func (command *EnableSonarCleanupCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "enable-sonar-cleanup",
		Usage:  "Check all repositories and enable sonar cleanup feature if sonar is enabled on repository",
		Action: command.EnableSonarCleanupAction,
	}
}

// EnableSonarCleanupAction allow to turn on the sonar cleanup setting on all available repositories
func (command *EnableSonarCleanupCommand) EnableSonarCleanupAction(context *cli.Context) error {
	cache := command.Settings.GetFileCache()
	client, _ := command.Settings.GetAPIClient()

	for _, repository := range cache.Repositories {
		sonarSettings, _ := client.GetSonarSettings(repository.Project.Key, repository.Slug)
		if sonarSettings.Project.SonarEnabled == true && sonarSettings.Project.ProjectCleanupEnabled == false {
			sonarSettings.Project.ProjectCleanupEnabled = true
			err := client.SetSonarSettings(repository.Project.Key, repository.Slug, sonarSettings)

			if err != nil {
				return err
			}

			fmt.Printf("[OK] Updated sonar cleanup settings for repository %s/%s\n", repository.Project.Key, repository.Slug)
		}
	}

	return nil
}
