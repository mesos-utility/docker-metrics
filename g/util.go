package g

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

// FileExists returns whether a file exists at a given filesystem path.
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("directory found instead of file at: %s", path)
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("not found: %s", path)
	}
	return false, err
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

		ret, err = FileExists(filepath.Join(dir, file))
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
