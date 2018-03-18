package gofaker

import (
	"path/filepath"
	"runtime"
	"strings"
)

var (
	gopath    string
	gopathlen int
)

func init() {
	pc, file, _, ok := runtime.Caller(0)
	if !ok || file == "?" {
		return
	}
	fn := runtime.FuncForPC(pc)
	fnStart := strings.LastIndex(fn.Name(), ".")
	if fnStart < 0 {
		return
	}
	fnPkg := fn.Name()[:strings.LastIndex(fn.Name(), "gofaker.init")]
	fnPkgStart := strings.Index(file, fnPkg)
	if fnPkgStart < 0 {
		return
	}
	gopathlen = fnPkgStart
	gopath = file[:gopathlen]
}

func trimGOPATH(filename string) string {
	if strings.HasPrefix(filename, gopath) {
		return filename[gopathlen:]
	}
	return filepath.Base(filename)
}

func GetSourceCodeLine(skip int) (file string, line int) {
	pc := make([]uintptr, 7)
	n := runtime.Callers(skip, pc)
	if n > 0 {
		stack := runtime.CallersFrames(pc[:n])
		f, _ := stack.Next()
		file = trimGOPATH(f.File)
		line = f.Line
	} else {
		file = "unknown"
	}
	return
}
