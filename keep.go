package main

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"
)

func keepDir(dir string, regex *regexp.Regexp, keepLast int) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var modTime time.Time
	var names []string

	for _, fi := range files {
		if regex.MatchString(fi.Name()) {
			if fi.Mode().IsRegular() {
				if !fi.ModTime().Before(modTime) {
					if fi.ModTime().After(modTime) {
						modTime = fi.ModTime()
						names = names[:0]
					}
					names = append(names, fi.Name())
				}
			}
		}
	}

	if len(names) > keepLast {
		toRemove := names[keepLast:]

		for _, f := range toRemove {
			err = os.Remove(path.Join(dir, f))

			if err != nil {
				return err
			}
		}
	}

	return nil
}
