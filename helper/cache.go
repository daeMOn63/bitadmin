// Package helper provides handy func and struct to be reused in commands
package helper

import (
	"encoding/json"
	"fmt"
	"github.com/daeMOn63/bitclient"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Cache interface define contract for Cache implementations
type Cache interface {
	WriteObject(namespace string, object interface{}) error
	Clear(namespace string) error
}

// FileCache is a Cache implementation storing data in a file
type FileCache struct {
	cacheDir     string
	Users        []bitclient.User
	Projects     []bitclient.Project
	Repositories []bitclient.Repository
}

// SearchRepositorySlug lookup for given repository slug in cached repositories
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

// Save write the cached data to the file
func (c *FileCache) Save() error {
	data, _ := json.Marshal(c)
	filename := c.getCacheFileName()

	os.Mkdir(filepath.Dir(filename), 0775)

	err := ioutil.WriteFile(filename, data, 0644)

	return err
}

// Clear erase cached data both in memory and file
func (c *FileCache) Clear() error {

	c.Users = nil
	c.Projects = nil
	c.Repositories = nil

	return c.Save()
}

// Load read file data and set them in memory
func (c *FileCache) Load() error {
	data, err := ioutil.ReadFile(c.getCacheFileName())
	if err == nil {
		json.Unmarshal(data, c)
	}

	return err
}

// String convert cached data to printable strings
func (c *FileCache) String() string {
	output := ""
	for _, user := range c.Users {
		output += fmt.Sprintf("user %d - %s - %s - %s - %s\n", user.Id, user.EmailAddress, user.Name, user.DisplayName, user.Slug)
	}

	for _, project := range c.Projects {

		var projectLinks []string
		for _, sublinks := range project.Links {
			for _, link := range sublinks {
				projectLinks = append(projectLinks, link["href"])
			}
		}

		output += fmt.Sprintf("project %s - %s - %s\n", project.Key, project.Name, strings.Join(projectLinks, " - "))
	}

	for _, repo := range c.Repositories {
		var repoLinks []string
		for _, sublinks := range repo.Links {
			for _, link := range sublinks {
				repoLinks = append(repoLinks, link["href"])
			}
		}

		output += fmt.Sprintf("repository %s/%s - %s - %s\n", repo.Project.Key, repo.Slug, repo.Name, strings.Join(repoLinks, " - "))
	}

	return output
}

// NewFileCache create a new FileCache instance
func NewFileCache(cacheDir string) *FileCache {
	return &FileCache{
		cacheDir: cacheDir,
	}
}
