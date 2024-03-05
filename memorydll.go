package memorydll

/*
#include "MemoryModule.h"
*/
import "C"
import (
	"errors"
	"syscall"
	"unsafe"
)

type Handle = C.HMEMORYMODULE

// A DLL implements access to a single DLL.
type DLL struct {
	Name   string
	Handle Handle
}

// DLLError describes reasons for DLL load failures.
type DLLError struct {
	Err     error
	ObjName string
	Msg     string
}

func (e *DLLError) Error() string { return e.Msg }

// FindProc searches DLL d for procedure named name and returns *Proc
// if found. It returns an error if search fails.
func (d *DLL) FindProc(name string) (proc *Proc, err error) {
	return memoryGetProcAddress(d, name)
}

// MustFindProc is like FindProc but panics if search fails.
func (d *DLL) MustFindProc(name string) *Proc {
	p, e := d.FindProc(name)
	if e != nil {
		panic(e)
	}
	return p
}

// Release unloads DLL d from memory.
func (d *DLL) Release() {
	memoryFreeLibrary(d)
}

// A Proc implements access to a procedure inside a DLL.
type Proc struct {
	Dll  *DLL
	Name string
	addr uintptr
}

// Addr returns the address of the procedure represented by p.
// The return value can be passed to Syscall to run the procedure.
func (p *Proc) Addr() uintptr {
	return p.addr
}

// Call executes procedure p with arguments a. It will panic, if more then 15 arguments
// are supplied.
//
// The returned error is always non-nil, constructed from the result of GetLastError.
// Callers must inspect the primary return value to decide whether an error occurred
// (according to the semantics of the specific function being called) before consulting
// the error. The error will be guaranteed to contain syscall.Errno.
func (p *Proc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	return syscall.SyscallN(p.Addr(), a...)
}

// remember to release when this dll is useless
func NewDLL(data []byte, name string) (*DLL, error) {
	ptr := unsafe.Pointer(&data[0])
	handle := C.MemoryLoadLibrary(ptr, C.size_t(len(data)))
	if handle != nil {
		return &DLL{
			Name:   name,
			Handle: handle,
		}, nil
	} else {
		e := errors.New("dll data error")
		return nil, &DLLError{
			Err:     e,
			ObjName: name,
			Msg:     "Failed to load " + name + ": " + e.Error(),
		}
	}

}

func memoryGetProcAddress(dll *DLL, procName string) (proc *Proc, err error) {
	cname := C.CString(procName)
	defer C.free(unsafe.Pointer(cname))

	addr := C.MemoryGetProcAddress(dll.Handle, cname)
	if addr != nil {
		return &Proc{
			Dll:  dll,
			Name: procName,
			addr: uintptr(unsafe.Pointer(addr)),
		}, nil
	}
	e := errors.New("no such function")
	return nil, &DLLError{
		Err:     e,
		ObjName: procName,
		Msg:     "Failed to find " + procName + " procedure in " + dll.Name + ": " + e.Error(),
	}
}

// remember free!
func memoryFreeLibrary(dll *DLL) {
	C.MemoryFreeLibrary(dll.Handle)
}
