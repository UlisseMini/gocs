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
	homedir "github.com/mitchellh/go-homedir"
)

// Config file struct, will be stored in ~/.goc.yaml
type Config struct {
	Github string // Github username
	Author string // Full name of the author
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	conf := getConfig()

	log.Info(conf)
}

// getConfig gets the config file from ~/.goc.yaml and returns it,
// if it does not exist it creates it.
func getConfig() (conf Config) {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
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
