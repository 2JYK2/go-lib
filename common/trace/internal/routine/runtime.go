package routine

import (
	"reflect"
	"unsafe"

	_ "github.com/2JYK2/go-lib/common/trace/internal/routine/g"
)

// getgp returns the pointer to the current runtime.g.
//
//go:linkname getgp fee-api-idc/common/trace/internal/routine/g.getgp
func getgp() unsafe.Pointer

// getgt returns the type of runtime.g.
//
//go:linkname getgt fee-api-idc/common/trace/internal/routine/g.getgt
func getgt() reflect.Type
