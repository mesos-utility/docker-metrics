package g

import (
	"runtime"
)

const (
	IDLEN int = 12
)

var VERSION = ""

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
