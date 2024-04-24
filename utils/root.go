package utils

import (
	"runtime"

	"github.com/gookit/goutil/fsutil"
)

func ProjectRoot() string {
	_, fullfilename, _, _ := runtime.Caller(0)
	return fsutil.DirPath(fullfilename + "/../../")
}
