package main

/*
#include <stdlib.h>
#cgo pkg-config: libevdev
#include <libevdev/libevdev.h>
#include <libevdev/libevdev-uinput.h>
*/
import "C"

import (
	"fmt"
	"os"
	"time"
)

func main() {
	vDev := NewVirtualMouse("Mouser")
	defer vDev.Close()
	fmt.Println(vDev.Move(10, 10))
}

type VirtualMouse struct {
	uidev *C.struct_libevdev_uinput
	dev   *C.struct_libevdev
}

func NewVirtualMouse(name string) *VirtualMouse {
	var dev *C.struct_libevdev
	var uidev *C.struct_libevdev_uinput

	dev = C.libevdev_new()

	C.libevdev_set_name(dev, C.CString(name))

	C.libevdev_enable_event_type(dev, C.EV_REL)
	C.libevdev_enable_event_code(dev, C.EV_REL, C.REL_X, nil)
	C.libevdev_enable_event_code(dev, C.EV_REL, C.REL_Y, nil)
	C.libevdev_enable_event_type(dev, C.EV_KEY)
	C.libevdev_enable_event_code(dev, C.EV_KEY, C.BTN_LEFT, nil)
	C.libevdev_enable_event_code(dev, C.EV_KEY, C.BTN_MIDDLE, nil)
	C.libevdev_enable_event_code(dev, C.EV_KEY, C.BTN_RIGHT, nil)

	rv := C.libevdev_uinput_create_from_device(dev, C.LIBEVDEV_UINPUT_OPEN_MANAGED, &uidev)
	if rv > 0 {
		fmt.Fprintf(os.Stderr, "Failed to create new uinput device: %v", rv)
		return nil
	}
	fmt.Printf("Virtual keyboard created at %s\n", C.GoString(C.libevdev_uinput_get_devnode(uidev)))
	time.Sleep(1 * time.Second)

	return &VirtualMouse{uidev: uidev, dev: dev}
}

func (u *VirtualMouse) Move(x int, y int) int {
	xD := C.libevdev_uinput_write_event(u.uidev, C.EV_REL, C.REL_X, C.int(x))
	yD := C.libevdev_uinput_write_event(u.uidev, C.EV_REL, C.REL_Y, C.int(y))
	if xD != 0 || yD != 0 {
		fmt.Printf("Mouse move failed x %d y %d", x, y)
		return -1
	}
	return 0
}

func (u *VirtualMouse) Close() {
	C.libevdev_uinput_destroy(u.uidev)
	C.libevdev_free(u.dev)
}
