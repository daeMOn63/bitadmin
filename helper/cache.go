package helper

import (
	"encoding/json"
	"fmt"
	"github.com/daeMOn63/bitclient"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Cache interface {
	WriteObject(namespace string, object interface{}) error
	Clear(namespace string) error
}

type FileCache struct {
	cacheDir     string
	Users        []bitclient.User
	Projects     []bitclient.Project
	Repositories []bitclient.Repository
}

func (c *FileCache) SearchRepositorySlug(slug string) (*bitclient.Repository, error) {
	for _, repo := range c.Repositories {
		if repo.Slug == slug {
			return &repo, nil
		}
	}

	return nil, fmt.Errorf("Cannot find repositry with slug {%s}, maybe refresh local cache ?", slug)
}

func (c *FileCache) getCacheFileName() string {
	return fmt.Sprintf("%s/cache", c.cacheDir)
}

func (c *FileCache) Save() error {
	data, _ := json.Marshal(c)
	filename := c.getCacheFileName()

	_ = os.Mkdir(filepath.Dir(filename), 0775)
	err := ioutil.WriteFile(filename, data, 0644)

	return err
}

func (c *FileCache) Clear() error {

	c.Users = nil
	c.Projects = nil
	c.Repositories = nil

	return c.Save()
}

func (c *FileCache) Load() error {
	data, err := ioutil.ReadFile(c.getCacheFileName())
	if err == nil {
		json.Unmarshal(data, c)
	}

	return err
}

func NewFileCache(cacheDir string) *FileCache {
	return &FileCache{
		cacheDir: cacheDir,
	}
}
