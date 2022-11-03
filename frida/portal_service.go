package frida

//#include <frida-core.h>
import "C"
import (
	"unsafe"
)

// Portal represents portal to collect exposed gadgets and sessions.
type Portal struct {
	portal *C.FridaPortalService
}

// NewPortal creates new Portal from the EndpointParameters provided.
func NewPortal(clusterParams, controlParams *EndpointParameters) *Portal {
	p := C.frida_portal_service_new(clusterParams.params, controlParams.params)

	return &Portal{
		portal: p,
	}
}

// Device returns portal device.
func (p *Portal) Device() *Device {
	dev := C.frida_portal_service_get_device(p.portal)
	return &Device{dev}
}

// ClusterParams returns the cluster parameters for the portal.
func (p *Portal) ClusterParams() *EndpointParameters {
	params := C.frida_portal_service_get_cluster_params(p.portal)
	return &EndpointParameters{params}
}

// ControlParams returns the control parameters for the portal.
func (p *Portal) ControlParams() *EndpointParameters {
	params := C.frida_portal_service_get_control_params(p.portal)
	return &EndpointParameters{params}
}

// Start stars the portal.
func (p *Portal) Start() error {
	var err *C.GError
	C.frida_portal_service_start_sync(p.portal, nil, &err)
	if err != nil {
		return &FError{err}
	}
	return nil
}

// Stop stops the portal.
func (p *Portal) Stop() error {
	var err *C.GError
	C.frida_portal_service_stop_sync(p.portal, nil, &err)
	if err != nil {
		return &FError{err}
	}
	return nil
}

// Kick kicks the connection with connectionID provided.
func (p *Portal) Kick(connectionID uint) {
	C.frida_portal_service_kick(p.portal, C.guint(connectionID))
}

// Post posts the message to the connectionID with json string or bytes.
func (p *Portal) Post(connectionID uint, json string, data []byte) {
	jsonC := C.CString(json)
	defer C.free(unsafe.Pointer(jsonC))

	gBytesData := goBytesToGBytes(data)

	C.frida_portal_service_post(p.portal, C.guint(connectionID), jsonC, gBytesData)
	clean(unsafe.Pointer(gBytesData), unrefGObject)
}

// Narrowcast sends the message to all controllers tagged with the tag.
func (p *Portal) Narrowcast(tag, json string, data []byte) {
	tagC := C.CString(tag)
	defer C.free(unsafe.Pointer(tagC))

	jsonC := C.CString(json)
	defer C.free(unsafe.Pointer(jsonC))

	gBytesData := goBytesToGBytes(data)
	C.frida_portal_service_narrowcast(p.portal, tagC, jsonC, gBytesData)
	clean(unsafe.Pointer(gBytesData), unrefGObject)
}

// Broadcast sends the message to all controllers.
func (p *Portal) Broadcast(json string, data []byte) {
	jsonC := C.CString(json)
	defer C.free(unsafe.Pointer(jsonC))

	gBytesData := goBytesToGBytes(data)
	C.frida_portal_service_broadcast(p.portal, jsonC, gBytesData)
	clean(unsafe.Pointer(gBytesData), unrefGObject)
}

// EnumerateTags returns all the tags that connection with connectionID is tagged
// with.
func (p *Portal) EnumerateTags(connectionID uint) []string {
	var length C.gint
	tagsC := C.frida_portal_service_enumerate_tags(
		p.portal,
		C.guint(connectionID),
		&length)

	return cArrayToStringSlice(tagsC, C.int(length))
}

// TagConnection tags the connection with connectionID with the tag provided.
func (p *Portal) TagConnection(connectionID uint, tag string) {
	tagC := C.CString(tag)
	defer C.free(unsafe.Pointer(tagC))

	C.frida_portal_service_tag(p.portal, C.guint(connectionID), tagC)
}

// UntagConnection untags the connection with connectionID with the tag provided.
func (p *Portal) UntagConnection(connectionID uint, tag string) {
	tagC := C.CString(tag)
	defer C.free(unsafe.Pointer(tagC))

	C.frida_portal_service_untag(p.portal, C.guint(connectionID), tagC)
}

func (p *Portal) On(sigName string, fn interface{}) {
	connectClosure(unsafe.Pointer(p.portal), sigName, fn)
}
