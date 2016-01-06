package g

import (
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

// display version info.
func HandleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}
}
