package caller

import (
	"runtime"
	"reflect"
)

// GetFuncName get function name
func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
