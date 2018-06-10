package jscoreWorker

/*
#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
#include "bridge.h"
typedef void (*closure)();
*/
import "C"
import (
	"reflect"
	"unsafe"
	"strings"
)

//export jsSend_go
func jsSend_go(ctx C.JSContextRef, channel C.int, value C.JSValueRef) C.JSValueRef {
	//println(channel)

	buf := NewJSUint8ArrayFromRef(ctx, NewJSValueFromRef(ctx, value).Object().ref).Array()
	if goListenersMap[int(channel)] != nil {
		goListenersMap[int(channel)](buf) // go call
	}
	return NewJSUndefinedFromCtxRef(ctx).ref
}
//export jsRecv_go
func jsRecv_go(ctx C.JSContextRef, channel C.int, listener C.JSObjectRef) C.JSValueRef {
	jsListenersMap[int(channel)] = append(jsListenersMap[int(channel)], listener)
	return NewJSUndefinedFromCtxRef(ctx).ref
}
//export jsPrint_go
func jsPrint_go(ctx C.JSContextRef, channel C.int, anyRef *C.JSObjectRef, anyCount int) C.JSValueRef {
	var anyList []C.JSObjectRef
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&anyList))
	sliceHeader.Cap = anyCount
	sliceHeader.Len = anyCount
	sliceHeader.Data = uintptr(unsafe.Pointer(anyRef))
	strList := make([]string, anyCount)
	for _, any := range anyList {
		strList = append(strList, NewJSObjectFromRef(ctx, any).Value().String())
	}
	println(strings.Join(strList, " "))
	return NewJSUndefinedFromCtxRef(ctx).ref
}

var jsListenersMap map[int][]C.JSObjectRef
var goListenersMap map[int]func(buf []byte)

func init() {
	jsListenersMap = make(map[int][]C.JSObjectRef)
	goListenersMap = make(map[int]func(buf []byte))
}

func GoSend(ctx C.JSContextRef, buf []byte, channel int) {
	args := make([]C.JSValueRef, 1)
	args[0] = NewJSUint8Array(ctx, buf).JSValue().ref
	jsErr := NewJSError(ctx)
	for _, listener := range jsListenersMap[channel] {
		C.JSObjectCallAsFunction(
			ctx,
			listener,
			NewJSObject(ctx).ref,
			1,
			&args[0],
			&jsErr.ref)
		if jsErr.ref != nil {
			panic("js err")
		}
	}
}

func GoRecv(ctx C.JSContextRef, listener func(buf []byte), channel int) {
	goListenersMap[channel] = listener
}