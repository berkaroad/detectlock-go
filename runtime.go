package detectlock

import (
	"bytes"
	"runtime"
	"strconv"
)

var getGoroutineID func() uint64

func init() {
	SetGetGoroutineIDFunc(func() uint64 {
		buf := make([]byte, 64)
		buf = buf[0:runtime.Stack(buf, false)]
		index := bytes.Index(buf, []byte{'['})
		buf = buf[0:index]
		buf = bytes.TrimLeft(buf, "goroutine")
		buf = bytes.TrimSpace(buf)
		gid, _ := strconv.ParseUint(string(buf), 10, 64)
		return gid
	})
}

// SetGetGoroutineIDFunc to set how to get goroutine id.
func SetGetGoroutineIDFunc(f func() uint64) {
	if f != nil {
		getGoroutineID = f
	}
}
