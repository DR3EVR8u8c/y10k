package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Yumfile struct {
	YumRepos        []YumRepoMirror `json:"repos"`
	LocalPathPrefix string          `json:"pathPrefix"`
}

var boolMap = map[bool]int{
	true:  1,
	false: 0,
}

// LoadYumfile loads a Yumfile from a json formated file
func LoadYumfile(path string) (*Yumfile, error) {
	Dprintf("Loading Yumfile: %s\n", path)

	yumfile := Yumfile{}

	// open file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// decode
	j := json.NewDecoder(f)
	if err = j.Decode(&yumfile); err != nil {
		return nil, err
	}

	// validate
	if err = yumfile.Validate(); err != nil {
		return nil, err
	}

	return &yumfile, nil
}

// Validate ensures all Yumfile fields contain valid values
func (c *Yumfile) Validate() error {
	for i, mirror := range c.YumRepos {
		if err := mirror.Validate(); err != nil {
			return err
		}

		// append path prefix
		if c.LocalPathPrefix != "" {
			c.YumRepos[i].LocalPath = fmt.Sprintf("%s/%s", c.LocalPathPrefix, mirror.LocalPath)
		}
	}

	return nil
}

func (c *Yumfile) Repo(id string) *YumRepoMirror {
	for _, mirror := range c.YumRepos {
		if mirror.YumRepo.ID == id {
			return &mirror
		}
	}

	return nil
}

// Sync processes all repository mirrors defined in a Yumfile
func (c *Yumfile) Sync() error {
	// sync each repo
	for _, mirror := range c.YumRepos {
		if err := mirror.Sync(); err != nil {
			return err
		}

		if err := mirror.Update(); err != nil {
			return err
		}
	}

	return nil
}
