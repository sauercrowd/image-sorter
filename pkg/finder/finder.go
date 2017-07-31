package finder

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

func Worker(ch chan string, inDir string, regex string) error {
	r := regexp.MustCompile(regex)
	err := find(ch, inDir, r)
	close(ch)
	return err
}

func find(ch chan string, inDir string, regex *regexp.Regexp) error {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		//if its a directory, recurse into it
		if f.IsDir() {
			find(ch, filepath.Join(inDir, f.Name()), regex)
			continue
		}
		if regex.FindString(f.Name()) == f.Name() {
			ch <- filepath.Join(inDir, f.Name())
		}
	}
	return nil
}
