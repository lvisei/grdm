package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lvisei/grdm"
)

var (
	// these are set in build step
	version = "unversioned"
	//lint:ignore U1000 embedded by goreleaser
	commit = "?"
	//lint:ignore U1000 embedded by goreleaser
	date = "?"

	downloadUrl   = flag.String("u", "", "Download URL")
	fileName      = flag.String("f", "", "Save file name with extension, defalut get from URL")
	dirPath       = flag.String("d", ".", "Save file Directory")
	totalSections = flag.Int("n", 2, "Number of connections to make to the server")
)

func main() {
	flag.Parse()

	if err := parse(); err == nil {
		startTime := time.Now()
		d := grdm.Download{
			Url:      *downloadUrl,
			FileName: *fileName,
			FileDir:  *dirPath,
			Sections: *totalSections,
		}
		filePath, err := d.Do()
		if err != nil {
			log.Printf("An error occured: %s\n", err)
		} else {
			fmt.Printf("Download completed in %v seconds\n", time.Since(startTime).Seconds())
			fmt.Printf("Saved to %s\n", filePath)

		}
	} else {
		fmt.Println(err)
		usage()
	}

}

func parse() error {
	if *downloadUrl == "" {
		return errors.New("download URL should not be empty")
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `
grdm

Version: v%s
HomePage: github.com/lvisei/grdm
Issue   : github.com/lvisei/grdm/issues
Author  : lvisei

Usage: grdm -u <URL> [-f] [-d] [-n]

Options:
	`, version)
	flag.PrintDefaults()
}
