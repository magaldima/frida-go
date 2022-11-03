package frida

//#include <frida-core.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Bus represent bus used to communicate with the devices.
type Bus struct {
	bus *C.FridaBus
}

// IsDetached returns whether the bus is detached from the device or not.
func (b *Bus) IsDetached() bool {
	dt := C.int(C.frida_bus_is_detached(b.bus))
	return dt == 1
}

// Attach attaches on the device bus.
func (b *Bus) Attach() error {
	var err *C.GError
	C.frida_bus_attach_sync(b.bus, nil, &err)
	if err != nil {
		return &FError{err}
	}
	return nil
}

// Post send(post) msg to the device.
func (b *Bus) Post(msg string, data []byte) {
	msgC := C.CString(msg)
	defer C.free(unsafe.Pointer(msgC))

	arr, sz := uint8ArrayFromByteSlice(data)
	defer C.free(unsafe.Pointer(arr))

	gBytesData := C.g_bytes_new((C.gconstpointer)(unsafe.Pointer(arr)), C.gsize(sz))
	C.frida_bus_post(b.bus, msgC, gBytesData)
	clean(unsafe.Pointer(gBytesData), unrefGObject)
}

// Clean will clean resources held by the bus.
func (b *Bus) Clean() {
	clean(unsafe.Pointer(b.bus), unrefFrida)
}

func (b *Bus) On(sigName string, fn interface{}) {
	connectClosure(unsafe.Pointer(b.bus), sigName, fn)
}
