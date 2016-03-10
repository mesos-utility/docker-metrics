package g

import (
	"os"
	"runtime"
	"testing"
)

func TestFileExists(t *testing.T) {
	_, file, _, _ := runtime.Caller(1)
	path, _ := os.Getwd()

	files := map[bool]string{
		true:  file,
		false: path,
	}

	for b, path := range files {
		ret, _ := FileExists(path)

		if b != ret {
			t.Fatalf("expected %v, got %v", b, ret)
		}
	}
}
