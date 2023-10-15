// all related utilities for handling slice
package sliceUtil

import (
	"fmt"
	"log"
	"runtime"
)

type slice struct {
	Items []interface{}
}

// (Handles string type) looks through items for that matches the needle
func (s *slice) InItems(needle interface{}) bool {
	if len(s.Items) <= 0 {
		return false
	}

	for _, item := range s.Items {
		if needle == item {
			return true
		}
	}

	return false
}

func getCallerInfo(level int) (string, int) {
	// Get information about the caller at depth 1 (the immediate caller)
	_, file, line, ok := runtime.Caller(level)
	if ok {
		return file, line
	} else {
		// this will most likely never be reached
		return "Unknown", 0
	}
}

// call this first to be able to handle any slice type
func Use(o ...interface{}) *slice {
	var anonSlice []interface{}

	for {
		var ok bool

		// handle string slice
		var strS []string
		strS, ok = o[0].([]string)
		if ok {
			anonSlice = make([]interface{}, len(strS))
			for k, v := range strS {
				anonSlice[k] = v
			}
			break
		}

		// handle int slice
		var intS []int
		intS, ok = o[0].([]int)
		if ok {
			anonSlice = make([]interface{}, len(intS))
			for k, v := range intS {
				anonSlice[k] = v
			}
			break
		}

		// handle float32 slice
		flt32S, ok := o[0].([]float32)
		if ok {
			anonSlice = make([]interface{}, len(flt32S))
			for k, v := range flt32S {
				anonSlice[k] = v
			}
			break
		}

		// handle float64 slice
		flt64S, ok := o[0].([]float64)
		if ok {
			anonSlice = make([]interface{}, len(flt64S))
			for k, v := range flt64S {
				anonSlice[k] = v
			}
			break
		}

		// handle bool slice
		boolS, ok := o[0].([]bool)
		if ok {
			anonSlice = make([]interface{}, len(boolS))
			for k, v := range boolS {
				anonSlice[k] = v
			}
			break
		}

		// handle interface slice
		mixedS, ok := o[0].([]interface{})
		if ok {
			anonSlice = make([]interface{}, len(mixedS))
			copy(anonSlice, mixedS)
			break
		}

		file, line := getCallerInfo(2)
		log.Fatal(fmt.Sprintf("Passed unsupported type to slice.Use() called in %s on line %d", file, line))
		break
	}

	interfaceSlice := make([]interface{}, len(anonSlice))
	copy(interfaceSlice, anonSlice)

	return &slice{
		Items: interfaceSlice,
	}
}
