package common_logger

import (
	"runtime"
	"strings"
)

func Trace(skip int) string {
	// 2 & 3 are file & line
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "?"
	}

	//Sanitize function name
	fn := runtime.FuncForPC(pc)
	functionName := fn.Name()
	idx := strings.LastIndexByte(functionName, '.')
	if idx != -1 {
		functionName = functionName[idx+1:]
	}
	return functionName
}
