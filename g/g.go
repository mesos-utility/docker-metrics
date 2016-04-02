package g

import (
	"runtime"
)

const (
	VERSION     = "0.2.0"
	IDLEN   int = 12
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
