package main

import (
	"path/filepath"
	"text/template"

	log "github.com/UlisseMini/leetlog"
)

// Project contains all the options needed for generating a new project.
type Project struct {
	Config // Config contains information stored in the configuration file.

	// name of the project, a directory with this name will be created.
	Project string
	Year    int
}

// Create creates the project to the current directory + the project name,
// templateDir is the template to use from ~/.goc/templates, and error is returned
// if it does not exist.
func (p Project) Create(templateDir string) error {
	log.Debugf("Create (templateDir): %q", templateDir)
	tmpl := template.New(p.Project)

	d := DirCopy{
		Tmpl: tmpl,
		Data: p,
	}

	src := filepath.Join(goc, templates, templateDir)
	return d.Copy(src, p.Project)
}
