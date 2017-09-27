package repository

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/settings"
	"github.com/urfave/cli"
)

type EnableSonarCleanupCommand struct {
	Settings *settings.BitAdminSettings
}

func (command *EnableSonarCleanupCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:   "enable-sonar-cleanup",
		Usage:  "Enable sonar cleanup feature if sonar is enabled on repository",
		Action: command.EnableSonarCleanupAction,
	}
}

func (command *EnableSonarCleanupCommand) EnableSonarCleanupAction(context *cli.Context) error {
	cache := command.Settings.GetFileCache()
	client, _ := command.Settings.GetApiClient()

	for _, repository := range cache.Repositories {
		sonarSettings, _ := client.GetSonarSettings(repository.Project.Key, repository.Slug)
		if sonarSettings.Project.SonarEnabled == true && sonarSettings.Project.ProjectCleanupEnabled == false {
			sonarSettings.Project.ProjectCleanupEnabled = true
			err := client.SetSonarSettings(repository.Project.Key, repository.Slug, sonarSettings)

			if err != nil {
				return err
			}

			fmt.Printf("Updated sonar cleanup settings for repository %s/%s\n", repository.Project.Key, repository.Slug)
		}
	}

	return nil
}
