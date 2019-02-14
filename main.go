// Copyright 2019 Ulisse Mini. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/UlisseMini/leetlog"
	"github.com/go-yaml/yaml"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"github.com/mitchellh/go-homedir"
)

// Only used on first program run, then unpacked to ~/.gocs
var box = packr.New("", "./gocs_default")

// Config file struct, will be stored in ~/.gocs/config.yaml
type Config struct {
	Github string // Github username
	Author string // Full name of the author
}

const (
	templates   = "templates" // for my sanity if i ever change this
	defaultTmpl = "default"   // default template in ~/.gocs/templates/
)

var (
	home       string // ~/
	gocs       string // ~/.gocs
	configPath string // ~/.gocs/config.yaml
)

func init() {
	var err error
	home, err = homedir.Dir()
	if err != nil {
		log.Fatalf("getting home directory: %v", err)
	}

	gocs = filepath.Join(home, ".gocs")
	configPath = filepath.Join(home, ".gocs/config.yaml")
}

func main() {
	// manage flags
	flag.Usage = usage

	d := flag.Bool("d", false, "print debug logs")
	flag.Parse()

	// switch over flags
	switch {
	case *d:
		log.DefaultLogger.Ldebug.SetOutput(os.Stderr)
	}

	// get the template dir to use
	templateDir := flag.Arg(1)
	if templateDir == "" {
		templateDir = defaultTmpl
	}

	createDir()         // create ~/.gocs if needed
	conf := getConfig() // read ~/.gocs/config.yaml and create it if needed

	proj := Project{
		Config:  conf,
		Year:    time.Now().Year(),
		Project: flag.Arg(0),
	}

	if err := proj.Create(templateDir); err != nil {
		log.Fatal(err)
	}
}

// create the ~/.gocs directory if it does not exist.
func createDir() {
	if _, err := os.Stat(gocs); err == nil {
		return
	}

	// create ~/.gocs
	if err := os.Mkdir(gocs, 0755); err != nil {
		log.Fatal(err)
	}

	// unpack box into ~/.gocs
	err := box.Walk(func(path string, file packd.File) error {
		log.Debugf("walk: %s", path)

		dst := filepath.Join(gocs, path)
		// if it has a parent directory create it.
		if err := createParents(dst); err != nil {
			return fmt.Errorf("walk: createParents: %v", err)
		}

		b, err := box.Find(path)
		if err != nil {
			log.Debugf("walk: find %q: %v", path, err)
			return err
		}

		return ioutil.WriteFile(dst, b, 0666)
	})
	if err != nil {
		log.Fatal(err)
	}
}

// getConfig gets the config file from ~/.gocs/config.yaml and returns it,
// if it does not exist it creates it.
func getConfig() (conf Config) {
	// Read the config file
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		// if it is not a path error then fatal
		if _, ok := err.(*os.PathError); !ok {
			log.Fatal(err)
		}
		// otherwise create a config
		conf = createConfig(configPath)
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s: <project> [template]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "<project> is the Project's name")
	fmt.Fprintln(os.Stderr, "[template] is looked for in ~/.gocs/")
	flag.PrintDefaults()
}

// create parent directories for path.
func createParents(path string) error {
	dir := filepath.Dir(path)
	if dir == "" {
		return nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("MkdirAll: %v", err)
	}
	return nil
}
