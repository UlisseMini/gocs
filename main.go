// Copyright 2019 Ulisse Mini. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO:
// Create template directory from packr box if it does not exist. (~/.goc)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/UlisseMini/leetlog"
	"github.com/go-yaml/yaml"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	homedir "github.com/mitchellh/go-homedir"
)

// Config file struct, will be stored in ~/.goc.yaml
type Config struct {
	Github string // Github username
	Author string // Full name of the author
}

var home string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	var err error
	home, err = homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
}

var box = packr.New("templates", "./default_templates")

func main() {
	flag.Parse()

	conf := getConfig()
	createDir()

	log.Info(conf)
}

// create the ~/.goc directory if it does not exist.
func createDir() {
	goc := filepath.Join(home, ".goc")
	if _, err := os.Stat(goc); err == nil {
		return
	}

	log.Infof("Create %q", goc)
	if err := mkchdir(goc); err != nil {
		log.Fatal(err)
	}

	err := box.Walk(func(path string, file packd.File) error {
		// if it has a parent directory create it.
		if d := filepath.Dir(path); d != "." {
			if err := os.Mkdir(d, 0755); err != nil {
				log.Debugf("mkdir %q: %v", d, err)
			}
		}

		log.Debugf("walk: %s", path)
		b, err := box.Find(path)
		if err != nil {
			log.Debugf("find %q: %v", path, err)
			return err
		}

		return ioutil.WriteFile(path, b, 0666)
	})

	log.Debugf("chdir %q: %v", "..", os.Chdir(".."))
	if err != nil {
		log.Fatal(err)
	}
}

// getConfig gets the config file from ~/.goc.yaml and returns it,
// if it does not exist it creates it.
func getConfig() (conf Config) {
	path := filepath.Join(home, ".goc.yaml")

	// Read the config file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		// if it is not a path error then fatal
		if _, ok := err.(*os.PathError); !ok {
			log.Fatal(err)
		}
		// otherwise create a config
		conf = createConfig(path)
	}

	if err := yaml.Unmarshal(b, &conf); err != nil {
		log.Fatal(err)
	}

	return conf
}

// question the user and create a config file, return the config when done.
func createConfig(path string) (conf Config) {
	s := bufio.NewScanner(os.Stdin)

	conf.Author = input(s, "Full name or alias: ")
	conf.Github = input(s, "Github username: ")

	data, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(path, data, 0666); err != nil {
		log.Fatal(err)
	}
	log.Infof("Written to %q", path)

	return conf
}

// input helper
func input(s *bufio.Scanner, prompt string) string {
	log.Print(prompt)
	if !s.Scan() {
		log.Fatal(s.Err())
	}

	return s.Text()
}

// make and chdir into a directory
func mkchdir(path string) error {
	// create the directory
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	// chdir into it
	return os.Chdir(path)
}
