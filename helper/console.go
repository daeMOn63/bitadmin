// Package helper provides handy func and struct to be reused in commands
package helper

import (
	"fmt"
	"github.com/daeMOn63/bitclient"
	"github.com/urfave/cli"
	"strings"
)

// PrintLinks output the bitclient.Links to stdout in a readable way
func PrintLinks(l bitclient.Links) {
	for typez, sublinks := range l {
		fmt.Printf("%s:\n", typez)
		for _, link := range sublinks {
			name := "web"
			if len(link["name"]) > 0 {
				name = link["name"]
			}
			fmt.Printf("\t%s - %s\n", name, link["href"])
		}
	}
}

// AppAutoComplete extends the default autocomplete provided by urfave/cli by filtering on flags when - or -- is typed
func AppAutoComplete(c *cli.Context) {
	flags := c.GlobalFlagNames()
	for _, flag := range flags {
		names := strings.Split(flag, ",")
		for _, n := range names {
			n = strings.TrimSpace(n)
			if len(n) > 0 && !c.GlobalIsSet(n) {
				s := "--"
				if len(n) == 1 {
					s = "-"
				}
				fmt.Printf("%s%s ", s, n)
			}
		}
	}

	cli.DefaultAppComplete(c)
}

// AutoComplete enhance the autocompletion by responding to project / user / username / repository... flags and printing
// available values from the cache.
func AutoComplete(c *cli.Context, cache *FileCache) {
	args := c.Parent().Args()
	lastArg := args[len(args)-1]

	switch lastArg {
	case "--project":
		for _, project := range cache.Projects {
			fmt.Fprintln(c.App.Writer, project.Key)
		}
	case "--user":
	case "--username":
		for _, user := range cache.Users {
			fmt.Fprintln(c.App.Writer, user.Slug)
		}
	case "--repository":
	case "--sourceRepository":
	case "--targetRepository":
		for _, repo := range cache.Repositories {
			fmt.Fprintln(c.App.Writer, repo.Slug)
		}
	case "--permission":
		fmt.Println("REPO_READ REPO_WRITE REPO_ADMIN")
	default:
		if lastArg[:2] == "--" {
			return
		}

		flags := c.Command.Flags
		for _, flag := range flags {
			name := flag.GetName()
			_, isStringSliceFlag := flag.(cli.StringSliceFlag)
			if !c.IsSet(name) || isStringSliceFlag {
				s := "--"
				if len(name) == 1 {
					s = "-"
				}
				fmt.Printf("%s%s ", s, name)
			}
		}
	}

}
