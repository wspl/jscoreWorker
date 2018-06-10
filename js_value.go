package jscoreWorker

/*
#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
#include <JavaScriptCore/JSValueRef.h>
#include <JavaScriptCore/JSObjectRef.h>
#include <JavaScriptCore/JSStringRef.h>
*/
import "C"

type JSValue struct {
	ref C.JSValueRef
	ctx C.JSContextRef
}

func NewJSValueFromRef(ctx C.JSContextRef, ref C.JSValueRef) *JSValue {
	jsVal := new(JSValue)
	jsVal.ctx = ctx
	jsVal.ref = ref
	return jsVal
}

func NewJSUndefined(ctx C.JSContextRef) *JSValue {
	jsVal := new(JSValue)
	jsVal.ctx = ctx
	jsVal.ref = C.JSValueMakeUndefined(ctx)
	return jsVal
}

func NewJSUndefinedFromCtxRef(ref C.JSContextRef) *JSValue {
	jsVal := new(JSValue)
	jsVal.ref = C.JSValueMakeUndefined(ref)
	return jsVal
}

func (v *JSValue) String() string {
	jsErr := NewJSError(v.ctx)
	jsStr := NewJSStringFromRef(C.JSValueToStringCopy(v.ctx, v.ref, &jsErr.ref))
	defer jsStr.Dispose()
	return jsStr.String()
}

func (v *JSValue) Object() *JSObject {
	jsErr := NewJSError(v.ctx)
	ret := C.JSValueToObject(v.ctx, v.ref, &jsErr.ref)
	if jsErr.ref != nil {
		panic("js err")
	}
	return NewJSObjectFromRef(v.ctx, ret)
}

func (v *JSValue) Float64() float64 {
	jsErr := NewJSError(v.ctx)
	ret := C.JSValueToNumber(v.ctx, v.ref, &jsErr.ref)
	if jsErr.ref != nil {
		panic("js err")
	}
	return float64(ret)
}

func (v *JSValue) Uint8() uint8 {
	jsErr := NewJSError(v.ctx)
	ret := C.JSValueToNumber(v.ctx, v.ref, &jsErr.ref)
	if jsErr.ref != nil {
		panic("js err")
	}
	return uint8(ret)
}