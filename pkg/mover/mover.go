package mover

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
)

func Worker(ch chan string, outDir string, wg *sync.WaitGroup, regex string, seperator string, useExif bool) {
	defer wg.Done()
	splitted := strings.Split(regex, seperator)
	var rp, rs *regexp.Regexp
	if len(splitted) > 0 {
		rp = regexp.MustCompile(splitted[0])
	}
	if len(splitted) > 1 {
		rs = regexp.MustCompile(splitted[1])
	}
	for fname := range ch {
		if useExif {
			sortFileExif(fname, outDir)
			continue
		}
		sortFile(fname, rs, rp, outDir)
	}
}

func sortFileExif(fname string, outDir string) {
	f, err := os.Open(fname)
	if err != nil {
		log.Printf("Error while opening file: %v", err)
		return
	}
	x, err := exif.Decode(f)
	if err != nil {
		//log.Printf("Error while deconding %s exif: %v", fname, err)
		return
	}
	d, err := x.DateTime()
	if err != nil {
		//log.Printf("Error while getting date from exif data: %v\n", err)
		return
	}
	date := d.Format("2006")
	outDir = path.Join(outDir, date)
	newPath := path.Join(outDir, filepath.Base(fname))
	log.Printf("Move from %s to %s", fname, newPath)
	moveFile(outDir, newPath, fname)
}

func sortFile(fname string, rs *regexp.Regexp, rp *regexp.Regexp, outDir string) {
	name := filepath.Base(fname)
	if rs != nil {
		name = rs.ReplaceAllString(name, "")
	}
	if rp != nil {
		name = rp.ReplaceAllString(name, "")
	}
	outDir = filepath.Join(outDir, name)
	newPath := path.Join(outDir, filepath.Base(fname))
	log.Printf("File: %s in dir %s", fname, newPath)
	moveFile(outDir, newPath, fname)
}

func createDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func moveFile(outDir, newPath, oldPath string) {
	if err := createDirIfNotExist(outDir); err != nil {
		log.Printf("Could not create dir %s: %v", outDir, err)
		return
	}
	err := os.Rename(oldPath, newPath)
	if err == nil {
		return
	}
	potentialError := fmt.Sprintf("rename %s %s: invalid cross-device link", oldPath, newPath)
	if err.Error() != potentialError {
		log.Printf("Could not move %s to %s: %v", oldPath, newPath, err)
		return
	}
	data, err := ioutil.ReadFile(oldPath)
	if err != nil {
		log.Printf("Could not read data to move it: %v", err)
		return
	}
	// Write data to dst
	err = ioutil.WriteFile(newPath, data, 0644)
	if err != nil {
		log.Printf("Could not write data to move it: %v", err)
		return
	}
	if err := os.Remove(oldPath); err != nil {
		log.Printf("Could not remove old file: %v", err)
	}
}
