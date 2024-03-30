package grdm

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

type Download struct {
	// The URL to download
	Url string
	// The target file save to name with extension, optional, defalut get from url
	FileName string
	// The target file save to directory
	FileDir string
	// The target file save to path, optional, equal fmt.Sprintf("%v/%v", d.FileDir, d.FileName)
	FilePath string
	// Number of sections/connections to make to the server
	Sections int
}

// Start the download
func (d Download) Do() (string, error) {
	filePathWithExt, err := d.parseFilePath()
	if err != nil {
		return "", err
	}

	d.FilePath = filePathWithExt
	err = os.MkdirAll(filepath.Dir(d.FilePath), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Println("Checking URL")
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got status %v\n", resp.StatusCode)

	if resp.StatusCode > 299 {
		return "", fmt.Errorf("can't process, response is %v", resp.StatusCode)
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return "", err
	}
	fmt.Printf("Size is %v bytes\n", size)

	var sections = make([][2]int, d.Sections)
	eachSize := size / d.Sections
	fmt.Printf("Each size is %v bytes\n", eachSize)

	// example: if file size is 100 bytes, our section should like:
	// [[0 10] [11 21] [22 32] [33 43] [44 54] [55 65] [66 76] [77 87] [88 98] [99 99]]
	for i := range sections {
		if i == 0 {
			// starting byte of first section
			sections[i][0] = 0
		} else {
			// starting byte of other sections
			sections[i][0] = sections[i-1][1] + 1
		}

		if i < d.Sections-1 {
			// ending byte of other sections
			sections[i][1] = sections[i][0] + eachSize
		} else {
			// ending byte of other sections
			sections[i][1] = size - 1
		}
	}

	var wg sync.WaitGroup
	// download each section concurrently
	for i, s := range sections {
		wg.Add(1)
		go func(i int, s [2]int) {
			defer wg.Done()
			err = d.downloadSection(i, s)
			if err != nil {
				panic(err)
			}
		}(i, s)
	}
	wg.Wait()

	return d.FilePath, d.mergeFiles(sections)
}

// Parse file save to path
func (d Download) parseFilePath() (string, error) {
	if d.Url == "" {
		return "", errors.New("download URL should not be empty")
	}

	// parse URI valid
	url, err := url.ParseRequestURI(d.Url)
	if err != nil {
		return "", err
	}

	if d.FilePath != "" {
		return d.FilePath, nil
	}

	if d.FileName == "" {
		filename := path.Base(url.Path)
		if filename == "" {
			return "", errors.New("file name should not be empty")
		} else {
			d.FileName = filename
		}
	}

	filePath := fmt.Sprintf("%v/%v", d.FileDir, d.FileName)

	return filePath, nil
}

// Download a single section and save content to a tmp file
func (d Download) downloadSection(i int, c [2]int) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c[0], c[1]))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("can't process, response is %v", resp.StatusCode)
	}
	fmt.Printf("Downloaded %v bytes for section %v\n", resp.Header.Get("Content-Length"), i)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%v/section-%v.tmp", d.FileDir, i), b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Get a new http request
func (d Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.Url,
		nil,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Chrome")
	return r, nil
}

// Merge tmp files to single file and delete tmp files
func (d Download) mergeFiles(sections [][2]int) error {
	f, err := os.OpenFile(d.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	for i := range sections {
		tmpFileName := fmt.Sprintf("%v/section-%v.tmp", d.FileDir, i)
		b, err := os.ReadFile(tmpFileName)
		if err != nil {
			return err
		}
		n, err := f.Write(b)
		if err != nil {
			return err
		}
		err = os.Remove(tmpFileName)
		if err != nil {
			return err
		}
		fmt.Printf("%v bytes merged\n", n)
	}
	return nil
}
