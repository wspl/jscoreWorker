package jscoreWorker

/*
#include <stdlib.h>
#include <JavaScriptCore/JSBase.h>
#include <JavaScriptCore/JSContextRef.h>
#include "bridge.h"
typedef void (*closure)();
*/
import "C"

//export jsSend_go
func jsSend_go(ctx C.JSContextRef, value C.JSValueRef, channel C.int) C.JSValueRef {
	//println(channel)

	buf := NewJSUint8ArrayFromRef(ctx, NewJSValueFromRef(ctx, value).Object().ref).Array()
	for _, listener := range goListenersMap[int(channel)] {
		listener(buf)
	}
	return NewJSUndefinedFromCtxRef(ctx).ref
}
//export jsRecv_go
func jsRecv_go(ctx C.JSContextRef, listener C.JSObjectRef, channel C.int) C.JSValueRef {
	jsListenersMap[int(channel)] = append(jsListenersMap[int(channel)], listener)
	return NewJSUndefinedFromCtxRef(ctx).ref
}
//export jsPrint_go
func jsPrint_go(ctx C.JSContextRef, any C.JSObjectRef, channel C.int) C.JSValueRef {
	println(NewJSObjectFromRef(ctx, any).Value().String())
	return NewJSUndefinedFromCtxRef(ctx).ref
}

var jsListenersMap map[int][]C.JSObjectRef
var goListenersMap map[int][]func(buf []byte)

func init() {
	jsListenersMap = make(map[int][]C.JSObjectRef)
	goListenersMap = make(map[int][]func(buf []byte))
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
	goListenersMap[channel] = append(goListenersMap[channel], listener)
}
//
//func JSCore() {
//	ctx := NewJSContext()
//	global := ctx.GetGlobal()
//
//	rawSendFn := C.JSObjectMakeFunctionWithCallback(ctx.ref, NewJSString("send").ref, C.closure(C.jsSend))
//	rawRecvFn := C.JSObjectMakeFunctionWithCallback(ctx.ref, NewJSString("recv").ref, C.closure(C.jsRecv))
//
//	rawWorker := NewJSObject(ctx.Convert())
//	rawWorker.SetProperty("send", NewJSValueFromRef(ctx.Convert(), C.JSValueRef(rawSendFn)))
//	rawWorker.SetProperty("recv", NewJSValueFromRef(ctx.Convert(), C.JSValueRef(rawRecvFn)))
//	global.SetProperty("rawWorker", rawWorker.Value())
//
//	ctx.EvaluateScript(`
//const JSCoreWorker = {
//	send (buf) { rawWorker.send(buf, 1) },
//	recv (cb) { rawWorker.recv(cb, 1) }
//}
//	`, "__jscoreWorker.js")
//
//	GoRecv(ctx.Convert(), func(buf []byte) {
//		fmt.Println(buf)
//	}, 1)
//
//	ctx.EvaluateScript(`
//	JSCoreWorker.recv((buf) => JSCoreWorker.send(buf));
//`, "f.js").String()
//
//	GoSend(ctx.Convert(), []byte{1,2,3,4,5}, 1)
//}