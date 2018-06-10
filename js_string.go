package jscoreWorker

/*
#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
#include <JavaScriptCore/JSStringRef.h>
*/
import "C"
import (
	"unsafe"
	"syscall"
)

type JSString struct {
	ref C.JSStringRef
}

func NewJSString(str string) *JSString {
	var jsStr = new(JSString)
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	jsStr.ref = C.JSStringCreateWithUTF8CString(cStr)
	return jsStr
}

func NewJSStringFromRef(ref C.JSStringRef) *JSString {
	var jsStr = new(JSString)
	jsStr.ref = ref
	return jsStr
}

func (s *JSString) Dispose() {
	C.JSStringRelease(s.ref)
}

func (s *JSString) String() string {
	l := C.JSStringGetMaximumUTF8CStringSize((C.JSStringRef)(unsafe.Pointer(s.ref)))
	buffer := C.malloc(l)
	if buffer == nil {
		panic(syscall.ENOMEM)
	}
	defer C.free(buffer)
	C.JSStringGetUTF8CString((C.JSStringRef)(unsafe.Pointer(s.ref)), (*C.char)(buffer), l)

	ret := C.GoString((*C.char)(buffer))
	return ret
}

func (s *JSString) Value(ctx C.JSContextRef) *JSValue {
	ref := C.JSValueMakeString(ctx, s.ref)
	return NewJSValueFromRef(ctx, ref)
}