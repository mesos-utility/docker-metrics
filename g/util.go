package g

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

// check file exists or not.
func IsExists(file string) (ret bool, err error) {
	if _, err := os.Stat(file); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// check cert files exists and read ok.
func CheckFilesExist(dir string, files []string) (ret bool, err error) {
	if dir == "" || len(files) <= 0 {
		return false, fmt.Errorf("dir or files is nil")
	}

	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			return false, fmt.Errorf("have a nil file name")
		}

		ret, err = IsExists(filepath.Join(dir, file))
		if err != nil {
			return ret, err
		}
	}

	return true, nil
}

// display version info.
func HandleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}
}

func Assert(err error) {
	if err != nil {
		glog.Fatalf("Assert failed: %s", err)
	}
}
