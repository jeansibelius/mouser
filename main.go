package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	str := "Hello world"
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	C.fputs(cs, (*C.FILE)(C.stdout))

	fmt.Println("Gopher here")
}
