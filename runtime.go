package detectlock

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
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

func getCaller(minimumCallerDepth int) *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, 10)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, more := frames.Next(); more; f, more = frames.Next() {
		// If the caller isn't part of this package, we're done
		if !strings.Contains(f.Function, "github.com/berkaroad/detectlock-go.") {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}
