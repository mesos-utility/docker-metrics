package g

import (
	"errors"
	"fmt"
	"os"
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
		return false, errors.New("dir or files is nil")
	}

	for _, file := range files {
		ret, err = IsExists(fmt.Sprintf("%s/%s", dir, file))

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
