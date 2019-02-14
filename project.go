package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	log "github.com/UlisseMini/leetlog"
	"github.com/fatih/color"
)

// Project contains all the options needed for generating a new project.
type Project struct {
	Config // Config contains information stored in the configuration file.

	// name of the project, a directory with this name will be created.
	Project string
	Year    int
}

// Create creates the project to the current directory + the project name,
// templateDir is the template to use from ~/.gocs/templates, and error is returned
// if it does not exist.
func (p Project) Create(templateDir string) error {
	initExists := false // true if '__init__` exists and is executeable
	// directory to copy templates from
	src := filepath.Join(gocs, templates, templateDir)

	// check for executeable '__init__'
	initBinary := filepath.Join(src, "__init__")
	info, _ := os.Stat(initBinary)
	if info != nil && info.Mode()&64 != 0 {
		initExists = true
	}

	log.Debugf("Create tmpl (templateDir): %q", templateDir)
	tmpl := template.New(p.Project)
	d := DirCopy{
		Tmpl: tmpl,
		Data: p,

		// ignore '__init__' if it exists and is executeable.
		// NOTE: Always ignoring '__init__' would be cleaner and lead to less code,
		// But i like the idea of it checking the executeable bit first.
		Ignore: func() []string {
			if initExists {
				return []string{"__init__"}
			} else {
				return nil
			}
		}(),
	}

	if err := d.Copy(src, p.Project); err != nil {
		return err
	}

	if initExists {
		color.HiRed("\n---------- BEGIN __init__ ----------")

		cmd := exec.Command(initBinary)
		cmd.Dir = p.Project

		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		color.HiRed("----------- END __init__ -----------")
		return err
	}

	return nil
}
