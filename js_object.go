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

type JSObject struct {
	ref C.JSObjectRef
	ctx C.JSContextRef
}

func NewJSObject(ctx C.JSContextRef) *JSObject {
	jsObj := new(JSObject)
	jsObj.ref = C.JSObjectMake(ctx, nil, nil)
	jsObj.ctx = ctx
	return jsObj
}

func NewJSObjectFromRef(ctx C.JSContextRef, ref C.JSObjectRef) *JSObject {
	jsObj := new(JSObject)
	jsObj.ref = ref
	jsObj.ctx = ctx
	return jsObj
}

func (o *JSObject) SetProperty(key string, value *JSValue) {
	propertyKey := NewJSString(key)
	defer propertyKey.Dispose()
	jsErr := NewJSError(o.ctx)
	C.JSObjectSetProperty(
		o.ctx,
		o.ref,
		propertyKey.ref,
		value.ref,
		(C.JSPropertyAttributes)(0),
		&jsErr.ref)
	if jsErr.ref != nil {
		panic("js err")
	}
}

func (o *JSObject) GetProperty(key string) *JSValue {
	jsErr := NewJSError(o.ctx)
	ret := C.JSObjectGetProperty(o.ctx, o.ref, NewJSString(key).ref, &jsErr.ref)
	if jsErr.ref != nil {
		panic("js err")
	}
	return NewJSValueFromRef(o.ctx, ret)
}

func (o *JSObject) Value() *JSValue {
	return NewJSValueFromRef(o.ctx, C.JSValueRef(o.ref))
}