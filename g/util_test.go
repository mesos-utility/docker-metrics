package g

import (
	"os"
	"runtime"
	"testing"
)

func TestFileExists(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	path, _ := os.Getwd()

	type fileTest struct {
		result   bool
		filename string
	}

	var files = []fileTest{
		{true, file},
		{false, path},
		{false, "."},
		{false, "../"},
	}

	for _, test := range files {
		ret, _ := FileExists(test.filename)

		if test.result != ret {
			t.Errorf("FileExists(%v) = %v, want %v", test.filename, ret, test.result)
		}
	}
}
