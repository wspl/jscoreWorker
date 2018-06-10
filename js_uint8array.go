package jscoreWorker

/*
#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
#include <JavaScriptCore/JSValueRef.h>
#include <JavaScriptCore/JSTypedArray.h>

extern void id(void *data, void *ctx) {}
typedef void (*closure)();
*/
import "C"
import (
	"unsafe"
	"strconv"
)

type JSUint8Array struct {
	ctx C.JSContextRef
	ref C.JSObjectRef
}

func NewJSUint8ArrayFromRef(ctx C.JSContextRef, ref C.JSObjectRef) *JSUint8Array {
	val := new(JSUint8Array)
	val.ctx = ctx
	val.ref = ref
	return val
}

func NewJSUint8Array(ctx C.JSContextRef, buf []byte) *JSUint8Array {
	errVal := NewJSError(ctx)

	val := new(JSUint8Array)
	val.ctx = ctx

	val.ref = C.JSObjectMakeTypedArrayWithBytesNoCopy(
		ctx,
		C.kJSTypedArrayTypeUint8Array,
		unsafe.Pointer(&buf[0]),
		C.size_t(len(buf)),
		C.closure(C.id),
		nil,
		&errVal.ref)
	if errVal.ref != nil {
		panic(errVal)
	}
	return val
}

func (a *JSUint8Array) JSValue() *JSValue {
	return NewJSValueFromRef(a.ctx, C.JSValueRef(a.ref))
}

func (a *JSUint8Array) Array() []byte {
	jsObj := NewJSObjectFromRef(a.ctx, a.ref)
	l := int(jsObj.GetProperty("byteLength").Float64())
	buf := make([]byte, l)
	for i := 0; i < l; i++ {
		buf[i] = jsObj.GetProperty(strconv.Itoa(i)).Uint8()
	}
	return buf
}