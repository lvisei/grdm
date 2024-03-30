package grdm_test

import (
	"os"
	"testing"

	"github.com/lvisei/grdm"
)

func TestDo(t *testing.T) {
	d := grdm.Download{
		Url:      "https://go.dev/dl/go1.22.1.src.tar.gz",
		FileDir:  "tmp",
		Sections: 2,
	}

	filePath, err := d.Do()
	if err != nil {
		t.Fatal("download", err)
	}

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf(" %s is not exist download file", filePath)
		}
	}
}
