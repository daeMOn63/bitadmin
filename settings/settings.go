// Package settings define the global application settings & flags
package settings

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

// BitAdminSettings hold the global flags values
type BitAdminSettings struct {
	Username     string
	Password     string
	PasswordFile string
	URL          string
	TempDir      string
}

// GetFlags provide the []cli.Flag needed by a cli.Command
func (bs *BitAdminSettings) GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "user",
			Usage:       "Authenticate on bitbucket with `<username>`",
			Destination: &bs.Username,
		},
		cli.StringFlag{
			Name:        "url",
			Usage:       "`<url>` of the bitbucket server",
			Destination: &bs.URL,
		},
		cli.StringFlag{
			Name:        "password",
			Usage:       "Read password from `<file>`",
			Destination: &bs.PasswordFile,
		},
	}
}

// GetAPIClient create a new instance of bitclient.BitClient initialized with the flag values
func (bs *BitAdminSettings) GetAPIClient() (*bitclient.BitClient, error) {

	// Load password from password file, checking for proper file permissions
	if bs.PasswordFile != "" {
		fileInfo, err := os.Stat(bs.PasswordFile)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Cannot read password file %s", bs.PasswordFile)
		}

		// Ensure proper permission on password file or named pipe if used
		mode := fileInfo.Mode() - (fileInfo.Mode() & os.ModeNamedPipe)
		if mode != 0600 {
			return nil, fmt.Errorf("Wrong permission on password file, please run \"chmod 600 %s\"", bs.PasswordFile)
		}

		passFromFile, err := ioutil.ReadFile(bs.PasswordFile)
		if err != nil {
			return nil, err
		}

		bs.Password = string(passFromFile)
	}

	if err := bs.Validate(); err != nil {
		return nil, err
	}

	return bitclient.NewBitClient(bs.URL, bs.Username, bs.Password), nil
}

// GetFileCache create a new instance of helper.FileCache and load the data from disk
func (bs *BitAdminSettings) GetFileCache() *helper.FileCache {
	cache := helper.NewFileCache(bs.TempDir)
	cache.Load()
	return cache
}

// Validate check for errors in global flag values
func (bs *BitAdminSettings) Validate() error {

	if bs.Username == "" {
		return fmt.Errorf("global flag --user is required")
	}

	if bs.Password == "" {
		return fmt.Errorf("global flag --password is required")
	}

	if bs.URL == "" {
		return fmt.Errorf("global flag --url is required")
	}

	return nil
}

// NewSettings create a new instance of BitAdminSettings
func NewSettings() *BitAdminSettings {
	return &BitAdminSettings{
		TempDir: os.TempDir() + "/bitadmin",
	}
}
