package common

import (
	"fmt"
	"path"
	"runtime"
)

func Log(params ...any) {
	caller := get_caller_info()
	params = append([]any{caller}, params...)
	fmt.Println(params...)
}

func get_caller_info() string {
	pc, filename, line, _ := runtime.Caller(2)
	funcname := "inline"
	if name := runtime.FuncForPC(pc).Name(); name != "" {
		funcname = name
	}
	filename = path.Base(filename)
	funcname = path.Base(funcname)
	return fmt.Sprintf("%s:%d [%s]", filename, line, funcname)
}
