package frida

//#include <frida-core.h>
import "C"

// PortalMembership type is used to join portal with session.
type PortalMembership struct {
	mem *C.FridaPortalMembership
}

// ID returns the ID of the membership
func (p *PortalMembership) ID() uint {
	return uint(C.frida_portal_membership_get_id(p.mem))
}

// Terminate terminates the session membership
func (p *PortalMembership) Terminate() error {
	var err *C.GError
	C.frida_portal_membership_terminate_sync(p.mem, nil, &err)
	if err != nil {
		return &FError{err}
	}
	return nil
}
