// used for copying the ~/.goc directory and executing templates on every file.
// taken from 'https://github.com/otiai10/copy' and modified to work with templates.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

// DirCopy is used to copy a directory, it will execute
// Templ on every file using 'data'
type DirCopy struct {
	Tmpl *template.Template // possible issue with this being reused
	Data interface{}
}

// Copy copies src to dest, doesn't matter if src is a directory or a file
func (d DirCopy) Copy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return d.dispatch(src, dest, info)
}

// dispatch dispatches copy-funcs according to the mode.
// Because this "dispatch" could be called recursively,
// "info" MUST be given here, NOT nil.
func (d DirCopy) dispatch(src, dst string, info os.FileInfo) error {
	// parse the name as a template
	t, err := d.Tmpl.Parse(dst)
	if err != nil {
		return fmt.Errorf("parse filename %q: %v", dst, err)
	}

	s := &strings.Builder{}
	if err := t.Execute(s, d.Data); err != nil {
		return fmt.Errorf("execute filename %q: %v", dst, err)
	}

	fmt.Printf("\t%s %s\n", color.GreenString("create"), s.String())
	dst = s.String()

	if info.Mode()&os.ModeSymlink != 0 {
		return d.lcopy(src, dst, info)
	}
	if info.IsDir() {
		return d.dcopy(src, dst, info)
	}
	return d.fcopy(src, dst, info)
}

// fcopy is for just a file,
// with considering existence of parent directory
// and file permission.
func (d DirCopy) fcopy(src, dst string, info os.FileInfo) error {

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(dst, info.Mode()); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	t, err := d.Tmpl.Parse(string(b))
	if err != nil {
		return err
	}
	return t.Execute(f, d.Data)
}

// dcopy is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
func (d DirCopy) dcopy(srcdir, destdir string, info os.FileInfo) error {

	if err := os.MkdirAll(destdir, info.Mode()); err != nil {
		return err
	}

	contents, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return err
	}

	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err := d.dispatch(cs, cd, content); err != nil {
			// If any error, exit immediately
			return err
		}
	}
	return nil
}

// lcopy is for a symlink,
// with just creating a new symlink by replicating src symlink.
func (d DirCopy) lcopy(src, dst string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dst)
}
