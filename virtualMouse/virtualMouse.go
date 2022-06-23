package virtualMouse

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
	fmt.Printf("Virtual mouse created at %s\n", C.GoString(C.libevdev_uinput_get_devnode(uidev)))
	time.Sleep(1 * time.Second)

	return &VirtualMouse{uidev: uidev, dev: dev}
}

func (vm *VirtualMouse) WriteEvent(eventType C.uint, eventCode C.uint, value C.int) int {
	// Write event to evdev
	return int(C.libevdev_uinput_write_event(vm.uidev, eventType, eventCode, value))
}

func (vm *VirtualMouse) TerminateEvent() {
	// Terminate event so kernel processes it
	// See: https://gitlab.freedesktop.org/libevdev/libevdev/-/blob/master/libevdev/libevdev-uinput.h#L238
	C.libevdev_uinput_write_event(vm.uidev, C.EV_SYN, C.SYN_REPORT, 0)
}

type InputEvent struct {
	eventType, eventCode C.uint
	value                C.int
}

func (vm *VirtualMouse) SendEvents(events []InputEvent) {
	for i := range events {
		vm.WriteEvent(events[i].eventType, events[i].eventCode, events[i].value)
	}
	vm.TerminateEvent()

	// If the event was a button press, release respective button
	for i := range events {
		if events[i].eventType == C.EV_KEY {
			C.libevdev_uinput_write_event(vm.uidev, events[i].eventType, events[i].eventCode, 0)
			vm.TerminateEvent()
		}
	}
}

func (vm *VirtualMouse) Move(x int, y int) {
	var events []InputEvent
	events = append(events, InputEvent{C.EV_REL, C.REL_X, C.int(x)})
	events = append(events, InputEvent{C.EV_REL, C.REL_Y, C.int(y)})
	vm.SendEvents(events)
}

func (vm *VirtualMouse) Click(button string) {
	var events []InputEvent
	switch button {
	case "left":
		events = append(events, InputEvent{C.EV_KEY, C.BTN_LEFT, 1})
	case "right":
		events = append(events, InputEvent{C.EV_KEY, C.BTN_RIGHT, 1})
	case "middle":
		events = append(events, InputEvent{C.EV_KEY, C.BTN_MIDDLE, 1})
	default:
		fmt.Println("Button press failed. No valid button provided, got:", button)
	}
	vm.SendEvents(events)
}

func (u *VirtualMouse) Close() {
	C.libevdev_uinput_destroy(u.uidev)
	C.libevdev_free(u.dev)
}
