package main

import (
	"fmt"
	"log"
	"time"

	"github.com/liuvigongzuoshi/grdm/internal/app"
)

func main() {
	startTime := time.Now()
	d := app.Download{
		// Provide the URL to download,
		Url: "http://vfx.mtime.cn/Video/2019/03/09/mp4/190309153658147087.mp4",
		// Provide the target file name with extension
		FileName: "sample.mp4",
		// Provide the target file Directory
		TargetDirPath: "tmp",
		// Number of sections/connections to make to the server
		TotalSections: 10,
	}
	err := d.Do()
	if err != nil {
		log.Printf("An error occured while downloading the file: %s\n", err)
	}
	fmt.Printf("Download completed in %v seconds\n", time.Now().Sub(startTime).Seconds())
}
