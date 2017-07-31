package main

import (
	"flag"
	"log"
	"strings"
	"sync"

	"github.com/sauercrowd/image-sorter/pkg/finder"
	"github.com/sauercrowd/image-sorter/pkg/mover"
)

func main() {
	f := parseFlags()
	ch := make(chan string)
	if f.threads < 2 {
		f.threads = 2
	}

	var wg sync.WaitGroup
	wg.Add(f.threads - 1)

	for i := 0; i < f.threads-1; i++ {
		go mover.Worker(ch, f.outDir, &wg, f.regex, f.placeholder, f.exif)
	}
	wildcard := ""
	for i := 0; i < len(f.placeholder); i++ {
		wildcard += "."
	}
	finderRegex := strings.Replace(f.regex, f.placeholder, wildcard, 1)
	err := finder.Worker(ch, f.inDir, finderRegex)
	if err != nil {
		log.Fatalf("Error while finding files: %v", err)
	}
	wg.Wait()
}

type flags struct {
	inDir, outDir, regex, placeholder string
	fileTypes                         []string
	threads                           int
	exif                              bool
}

func parseFlags() flags {
	in := flag.String("in", "", "Directory which should be sorted")
	out := flag.String("out", "", "Directory where sorted stuff should be moved")
	regx := flag.String("regex", ".*", "pattern of files which should be matched")
	sep := flag.String("placeholder", "XXXX", "The placeholder for the pattern")
	threads := flag.Int("n", 2, "number of threads to use, mininum 2")
	useExif := flag.Bool("exif", false, "if true, uses exif data to get year and move it into a folder according to its year")
	flag.Parse()
	return flags{inDir: *in, outDir: *out, regex: *regx, placeholder: *sep, threads: *threads, exif: *useExif}
}
