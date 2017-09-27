package settings

import (
	"fmt"
	"github.com/daeMOn63/bitadmin/helper"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
	"os"
)

type BitAdminSettings struct {
	Username string
	Password string
	Url      string
	TempDir  string
}

func (bs *BitAdminSettings) GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "user",
			Usage:       "Authenticate on bitbucket with `<username>`",
			Destination: &bs.Username,
		},
		cli.StringFlag{
			Name:        "password",
			Usage:       "Authenticate on bitbucket with `<password>`",
			Destination: &bs.Password,
		},
		cli.StringFlag{
			Name:        "url",
			Usage:       "`<url>` of the bitbucket server",
			Destination: &bs.Url,
		},
	}
}

func (bs *BitAdminSettings) GetApiClient() (*bitclient.BitClient, error) {
	if err := bs.Validate(); err != nil {
		return nil, err
	}

	return bitclient.NewBitClient(bs.Url, bs.Username, bs.Password), nil
}

func (bs *BitAdminSettings) GetFileCache() *helper.FileCache {
	cache := helper.NewFileCache(bs.TempDir)
	cache.Load()
	return cache
}

func (bs *BitAdminSettings) Validate() error {
	if bs.Username == "" {
		return fmt.Errorf("global flag --user is required")
	}

	if bs.Password == "" {
		return fmt.Errorf("global flag --password is required")
	}

	if bs.Url == "" {
		return fmt.Errorf("global flag --url is required")
	}

	return nil
}

func NewSettings() *BitAdminSettings {
	return &BitAdminSettings{
		TempDir: os.TempDir() + "/bitadmin",
	}
}
